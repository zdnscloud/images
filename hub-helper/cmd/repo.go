package cmd

import (
	"strings"
	"encoding/json"
	"fmt"

	"github.com/zdnscloud/images/hub-helper/client"

	"github.com/urfave/cli"
)

const defaultPageSize = 100

func RepoCommand() cli.Command {
	return cli.Command{
		Name:   "repo",
		Usage:  "repo command(list|info|search)",
		Subcommands: cli.Commands{
			cli.Command{
				Name: "list",
				Usage: "list all repos",
				Action: listRepo,
			},
			cli.Command{
				Name: "info",
				Usage: "print repo info",
				Action: infoRepo,
			},
			cli.Command{
				Name: "search",
				Usage: "search repo by input key",
				Action: searchRepo,
			},
		},
		Before: loadConfig,
	}
}

func ReposCommand() cli.Command {
	return cli.Command{
		Name:   "repos",
		Usage:  "list all repos",
		Action: listRepo,
		Before: loadConfig,
	}
}

func listRepo(ctx *cli.Context) error {
	repos, err := getRepos(globalClient)
	if err != nil {
		return err
	}

	fmt.Printf("total repo count %v\n", len(repos))
	for _, r := range repos {
		fmt.Println(r.Name)
	}
	return nil
}

func infoRepo(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return fmt.Errorf("must input repo name")
	}
	repo := ctx.Args().Get(0)

	repos, err := getRepos(globalClient)
	if err != nil {
		return err
	}

	for _, r := range repos {
		if r.Name == repo {
			jsonContent, err := json.MarshalIndent(&r, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonContent))
			return nil
		}
	}
	return nil
}

func searchRepo(ctx *cli.Context) error {
	if ctx.NArg() == 0 {
		return fmt.Errorf("must input repo name")
	}
	repo := ctx.Args().Get(0)

	repos, err := getRepos(globalClient)
	if err != nil {
		return err
	}

	for _, r := range repos {
		if strings.Contains(r.Name, repo) {
			fmt.Println(r.Name)
		}
	}
	return nil
}


func getRepos(c *client.Client) ([]client.Repo, error) {
	repos := []client.Repo{}
	page := 1
	for {
		pageRepos, err := c.GetRepos(page, defaultPageSize)
		if err != nil {
			return nil, fmt.Errorf("get repos for namespace '%s' from DockerHub error: %s", c.Namespce(), err.Error())
		}
		repos = append(repos, pageRepos.Repos...)

		if len(pageRepos.Next) == 0 {
			break
		}
		page++
	}
	return repos, nil
}
