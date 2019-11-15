package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/zdnscloud/cement/log"
)

func pullImage(cli *client.Client, images map[string]string) error {
	for _, image := range images {
		resp, err := cli.ImagePull(context.TODO(), image, types.ImagePullOptions{})
		if err != nil {
			return fmt.Errorf("Pull image %s failed %s", image, err.Error())
		}
		defer resp.Close()
		io.Copy(ioutil.Discard, resp)
		log.Infof("Pull image %s succeed", image)
	}
	return nil
}

func saveImage(cli *client.Client, images map[string]string, fileName string) error {
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
	var inputFile string
	var outFile string
	flag.StringVar(&inputFile, "i", "image.json", "image list json file,it's content should be a json map")
	flag.StringVar(&outFile, "o", "image.tar", "image tar file name")
	flag.Parse()

	var images map[string]string
	f, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Load input file %s failed %s", inputFile, err.Error())
	}

	if err := json.Unmarshal(f, &images); err != nil {
		log.Fatalf("Load input file %s content failed %s", inputFile, err.Error())
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Create docker client failed %s", err)
	}

	if err := pullImage(cli, images); err != nil {
		log.Fatalf("Pull image failed ", err.Error())
	}
	if err := saveImage(cli, images, outFile); err != nil {
		log.Fatalf("Save image %s failed %s", outFile, err.Error())
	}
	log.Infof("Save image %s succeed", outFile)
}
