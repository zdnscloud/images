package cmd

import (
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/zdnscloud/images/hub-helper/client"

	"github.com/urfave/cli"
)

var globalClient *client.Client

func ConfigCommand() cli.Command {
	return cli.Command{
		Name:   "config",
		Usage:  "config [flags]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "user",
				Value: "zdnscloud",
			},
			cli.StringFlag{
				Name: "password",
				Required: true,
			},
		},
		Action: writeConfig,
	}
}

func writeConfig(ctx *cli.Context) error {
	credential := &client.LoginCredential{
		User: ctx.String("user"),
		Password: ctx.String("password"),
	}

	content, err := json.MarshalIndent(credential, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config failed %s", err.Error())
	}

	file, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("get config file path failed %s", err.Error())
	}

	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		return fmt.Errorf("write config file failed %s", err.Error())
	}
	return nil
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".hub-helper.json"), nil
}

func loadConfig(ctx *cli.Context) error {
	file, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("get config file path failed %s", err.Error())
	}
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("load config failed %s", err.Error())
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("load config failed %s", err.Error())
	}

	result := &client.LoginCredential{}
	if err := json.Unmarshal(content, result); err != nil {
		return fmt.Errorf("unmarshal config json failed %s", err.Error())
	}
	
	hubClient, err := client.NewClient(result.User, result.Password)
	if err != nil {
		return err
	}
	
	globalClient = hubClient
	return nil
}

