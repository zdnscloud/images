package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/zdnscloud/cement/log"
)

func pullComponentImage(cli *client.Client, images map[string]string) error {
	for _, image := range images {
		resp, err := cli.ImagePull(context.TODO(), image, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("Pull image %s failed", image)
		}
		defer resp.Close()
		io.Copy(ioutil.Discard, resp)
		log.Infof("Pull image %s succeed", image)
	}
	return nil
}

func saveComponentImage(cli *client.Client, images map[string]string, fileName string) error {
	saveImages := []string{}
	for _, image := range images {
		saveImages = append(saveImages, image)
	}
	resp, err := cli.ImageSave(context.TODO(), saveImages)
	if err != nil {
		return err
	}
	defer resp.Close()
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 1024)
	for {
		_, err := resp.Read(buf)
		if err == io.EOF {
			break
		}
		f.Write(buf)
	}
	return nil
}

func main() {
	log.InitLogger(log.Info)
	defer log.CloseLogger()
	var version string
	flag.StringVar(&version, "version", "v2.0", "singlecloud version")
	flag.Parse()
	image, ok := images[version]
	if !ok {
		log.Fatalf("Not found version %s in image list", version)
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("New docker client failed %s", err)
	}

	for _, component := range image {
		log.Infof("Begining pull %s images", component.Name)
		if err := pullComponentImage(cli, component.Images); err != nil {
			log.Fatalf("Pull component %s images failed %s", component.Name, err)
		}
		fileName := component.Name + "-images.tar"
		if err := saveComponentImage(cli, component.Images, fileName); err != nil {
			log.Fatalf("Save component %s images %s failed %s", component.Name, fileName, err)
		}
		log.Infof("Save component %s images %s succeed", component.Name, fileName)
	}
	log.Infof("Finished package zcloud images")
}
