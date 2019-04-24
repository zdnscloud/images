package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	backupBaseDir = "/backup"
	backupRetries = 4
	ServerPort    = "2379"
)

var commonFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "endpoints",
		Usage: "Etcd endpoints",
		Value: "127.0.0.1:2379",
	},
	cli.BoolFlag{
		Name:   "debug",
		Usage:  "Verbose logging information for debugging purposes",
		EnvVar: "ZDNSCLOUD_DEBUG",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "Backup name to take once",
	},
	cli.StringFlag{
		Name:   "cacert",
		Usage:  "Etcd CA client certificate path",
		EnvVar: "ZDNSCLOUD_CACERT",
	},
	cli.StringFlag{
		Name:   "cert",
		Usage:  "Etcd client certificate path",
		EnvVar: "ETCD_CERT",
	},
	cli.StringFlag{
		Name:   "key",
		Usage:  "Etcd client key path",
		EnvVar: "ETCD_KEY",
	},
	cli.StringFlag{
		Name:   "local-endpoint",
		Usage:  "Local backup download endpoint",
		EnvVar: "LOCAL_ENDPOINT",
	},
}

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	err := os.Setenv("ETCDCTL_API", "3")
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "Etcd Wrapper"
	app.Usage = "Utility services for Etcd cluster backup"
	app.Commands = []cli.Command{
		RollingBackupCommand(),
	}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func RollingBackupCommand() cli.Command {

	snapshotFlags := []cli.Flag{
		cli.DurationFlag{
			Name:  "creation",
			Usage: "Create backups after this time interval in minutes",
			Value: 5 * time.Minute,
		},
		cli.DurationFlag{
			Name:  "retention",
			Usage: "Retain backups within this time interval in hours",
			Value: 24 * time.Hour,
		},
		cli.BoolFlag{
			Name:  "once",
			Usage: "Take backup only once",
		},
	}

	snapshotFlags = append(snapshotFlags, commonFlags...)

	return cli.Command{
		Name:  "etcd-backup",
		Usage: "Perform etcd backup tools",
		Subcommands: []cli.Command{
			{
				Name:   "save",
				Usage:  "Take snapshot on all etcd hosts and backup",
				Flags:  snapshotFlags,
				Action: RollingBackupAction,
			},
			{
				Name:   "download",
				Usage:  "Download specified snapshot local endpoint",
				Flags:  commonFlags,
				Action: DownloadBackupAction,
			},
		},
	}
}

func SetLoggingLevel(debug bool) {
	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func RollingBackupAction(c *cli.Context) error {
	SetLoggingLevel(c.Bool("debug"))

	creationPeriod := c.Duration("creation")
	retentionPeriod := c.Duration("retention")
	etcdCert := c.String("cert")
	etcdCACert := c.String("cacert")
	etcdKey := c.String("key")
	etcdEndpoints := c.String("endpoints")
	if len(etcdCert) == 0 || len(etcdCACert) == 0 || len(etcdKey) == 0 {
		log.WithFields(log.Fields{
			"etcdCert":   etcdCert,
			"etcdCACert": etcdCACert,
			"etcdKey":    etcdKey,
		}).Errorf("Failed to find etcd cert or key paths")
		return fmt.Errorf("Failed to find etcd cert or key paths")
	}
	log.WithFields(log.Fields{
		"creation":  creationPeriod,
		"retention": retentionPeriod,
	}).Info("Initializing Rolling Backups")

	if c.Bool("once") {
		backupName := c.String("name")
		if len(backupName) == 0 {
			backupName = fmt.Sprintf("%s_etcd", time.Now().Format(time.RFC3339))
		}
		if err := CreateBackup(backupName, etcdCACert, etcdCert, etcdKey, etcdEndpoints); err != nil {
			return err
		}
		prefix := getNamePrefix(backupName)
		// we only clean named backups if we have a retention period and a cluster name prefix
		if retentionPeriod != 0 && len(prefix) != 0 {
			if err := DeleteNamedBackups(retentionPeriod, prefix); err != nil {
				return err
			}
		}
		return nil
	}
	backupTicker := time.NewTicker(creationPeriod)
	for {
		select {
		case backupTime := <-backupTicker.C:
			backupName := fmt.Sprintf("%s_etcd", backupTime.Format(time.RFC3339))
			CreateBackup(backupName, etcdCACert, etcdCert, etcdKey, etcdEndpoints)
			DeleteBackups(backupTime, retentionPeriod)
		}
	}
}

func CreateBackup(backupName string, etcdCACert, etcdCert, etcdKey, endpoints string) error {
	failureInterval := 15 * time.Second
	backupDir := fmt.Sprintf("%s/%s", backupBaseDir, backupName)
	var err error
	for retries := 0; retries <= backupRetries; retries++ {
		if retries > 0 {
			time.Sleep(failureInterval)
		}
		// check if the cluster is healthy
		cmd := exec.Command("etcdctl",
			fmt.Sprintf("--endpoints=[%s]", endpoints),
			"--cacert="+etcdCACert,
			"--cert="+etcdCert,
			"--key="+etcdKey,
			"endpoint", "health")
		data, err := cmd.CombinedOutput()

		if strings.Contains(string(data), "unhealthy") {
			log.WithFields(log.Fields{
				"error": err,
				"data":  string(data),
			}).Warn("Checking member health failed from etcd member")
			continue
		}

		cmd = exec.Command("etcdctl",
			fmt.Sprintf("--endpoints=[%s]", endpoints),
			"--cacert="+etcdCACert,
			"--cert="+etcdCert,
			"--key="+etcdKey,
			"snapshot", "save", backupDir)

		startTime := time.Now()
		data, err = cmd.CombinedOutput()
		endTime := time.Now()

		if err != nil {
			log.WithFields(log.Fields{
				"attempt": retries + 1,
				"error":   err,
				"data":    string(data),
			}).Warn("Backup failed")
			continue
		}
		log.WithFields(log.Fields{
			"name":    backupName,
			"runtime": endTime.Sub(startTime),
		}).Info("Created backup")
	}
	return err
}

func DeleteBackups(backupTime time.Time, retentionPeriod time.Duration) {
	files, err := ioutil.ReadDir(backupBaseDir)
	if err != nil {
		log.WithFields(log.Fields{
			"dir":   backupBaseDir,
			"error": err,
		}).Warn("Can't read backup directory")
	}

	cutoffTime := backupTime.Add(retentionPeriod * -1)

	for _, file := range files {
		if file.IsDir() {
			log.WithFields(log.Fields{
				"name": file.Name(),
			}).Warn("Ignored directory, expecting file")
			continue
		}

		backupTime, err2 := time.Parse(time.RFC3339, strings.Split(file.Name(), "_")[0])
		if err2 != nil {
			log.WithFields(log.Fields{
				"name":  file.Name(),
				"error": err2,
			}).Warn("Couldn't parse backup")

		} else if backupTime.Before(cutoffTime) {
			_ = DeleteBackup(file)
		}
	}
}

func DeleteBackup(file os.FileInfo) error {
	toDelete := fmt.Sprintf("%s/%s", backupBaseDir, file.Name())

	cmd := exec.Command("rm", "-r", toDelete)

	startTime := time.Now()
	err2 := cmd.Run()
	endTime := time.Now()

	if err2 != nil {
		log.WithFields(log.Fields{
			"name":  file.Name(),
			"error": err2,
		}).Warn("Delete backup failed")
		return err2
	}
	log.WithFields(log.Fields{
		"name":    file.Name(),
		"runtime": endTime.Sub(startTime),
	}).Info("Deleted backup")
	return nil
}

func DownloadBackupAction(c *cli.Context) error {
	return DownloadLocalBackup(c)
}

func DownloadLocalBackup(c *cli.Context) error {
	snapshot := path.Base(c.String("name"))
	endpoint := c.String("local-endpoint")
	if snapshot == "." || snapshot == "/" {
		return fmt.Errorf("snapshot name is required")
	}
	if len(endpoint) == 0 {
		return fmt.Errorf("local-endpoint is required")
	}
	certs, err := getCertsFromCli(c)
	if err != nil {
		return err
	}
	tlsConfig, err := setupTLSConfig(certs, false)
	if err != nil {
		return err
	}
	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	snapshotFile, err := os.Create(fmt.Sprintf("%s/%s", backupBaseDir, snapshot))
	if err != nil {
		return err
	}
	defer snapshotFile.Close()
	log.Infof("Invoking downloading backup files: %s", snapshot)
	resp, err := client.Get(fmt.Sprintf("https://%s:%s/%s", endpoint, ServerPort, snapshot))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(snapshotFile, resp.Body); err != nil {
		return err
	}
	log.Infof("Successfully download %s from %s ", snapshot, endpoint)
	return nil
}

func DeleteNamedBackups(retentionPeriod time.Duration, prefix string) error {
	files, err := ioutil.ReadDir(backupBaseDir)
	if err != nil {
		log.WithFields(log.Fields{
			"dir":   backupBaseDir,
			"error": err,
		}).Warn("Can't read backup directory")
		return err
	}
	cutoffTime := time.Now().Add(retentionPeriod * -1)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), prefix) && file.ModTime().Before(cutoffTime) {
			if err = DeleteBackup(file); err != nil {
				return err
			}
		}
	}
	return nil
}

func getNamePrefix(name string) string {
	re := regexp.MustCompile("^c-[a-z0-9].*?-")
	m := re.FindStringSubmatch(name)
	if len(m) == 0 {
		return ""
	}
	return m[0]
}

func getCertsFromCli(c *cli.Context) (map[string]string, error) {
	caCert := c.String("cacert")
	cert := c.String("cert")
	key := c.String("key")
	if len(cert) == 0 || len(caCert) == 0 || len(key) == 0 {
		return nil, fmt.Errorf("cacert, cert and key are required")
	}

	return map[string]string{"cacert": caCert, "cert": cert, "key": key}, nil
}

func setupTLSConfig(certs map[string]string, isServer bool) (*tls.Config, error) {
	caCertPem, err := ioutil.ReadFile(certs["cacert"])
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCertPem)
	if isServer {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = certPool
		tlsConfig.MinVersion = tls.VersionTLS12
	} else { // client config
		x509Pair, err := tls.LoadX509KeyPair(certs["cert"], certs["key"])
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{x509Pair}
		tlsConfig.RootCAs = certPool
		// This is to avoid IP SAN errors.
		tlsConfig.InsecureSkipVerify = true
	}

	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}
