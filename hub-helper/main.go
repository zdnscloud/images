package main

import (
	"fmt"
	"os"
	"log"

	"github.com/zdnscloud/images/hub-helper/cmd"

	"github.com/urfave/cli"
)

var (
	version string
	build   string
)

func main() {
	app := cli.NewApp()
	app.Name = "hub-helper"
	app.Version = fmt.Sprintf("hub-helper version %s build at %s", version, build)
	app.Usage = "hub-helper command [args]"
	app.Description = "dockerhub cli tool"
	app.Author = "Zcloud"
	app.Email = "zcloud@zdns.cn"
	app.Commands = []cli.Command{
		cmd.ConfigCommand(),
		cmd.RepoCommand(),
		cmd.ReposCommand(),
		cmd.TagCommand(),
		cmd.TagsCommand(),
		cmd.DeleteCommand(),

	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
