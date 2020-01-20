package cmd

import (
	"strings"
	"encoding/json"
	"fmt"

	"github.com/zdnscloud/images/hub-helper/client"

	"github.com/urfave/cli"
)

func TagCommand() cli.Command {
	return cli.Command{
		Name:   "tag",
		Usage:  "tag command(list|info|search)",
		Subcommands: cli.Commands{
			cli.Command{
				Name: "list",
				Usage: "list repo tags",
				Action: listTag,
			},
			cli.Command{
				Name: "info",
				Usage: "print repo tag info",
				Action: infoTag,
			},
			cli.Command{
				Name: "search",
				Usage: "search repo tag by input key",
				Action: searchTag,
			},
			cli.Command{
				Name: "delete",
				Usage: "delete repo tag by input key",
				Action: deleteTag,
			},
		},
		Before: loadConfig,
	}
}

func TagsCommand() cli.Command {
	return cli.Command{
		Name:   "tags",
		Usage:  "list repo tags",
		Action: listTag,
		Before: loadConfig,
	}
}

func DeleteCommand() cli.Command {
	return cli.Command{
		Name:   "delete",
		Usage:  "delete repo tag",
		Action: deleteTag,
		Before: loadConfig,
	}
}


func listTag(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("must input repo name")
	}
	repo := ctx.Args().Get(0)

	tags, err := getRepoTags(globalClient, repo)
	if err != nil {
		return err
	}

	fmt.Printf("total tag count %v\n", len(tags))
	for _, t := range tags {
		fmt.Println(t.Name)
	}
	return nil
}

func infoTag(ctx *cli.Context) error {
	repo, tag, err := getTagArgs(ctx)
	if err != nil {
		return err
	}

	tags, err := getRepoTags(globalClient, repo)
	if err != nil {
		return err
	}

	for _, t := range tags {
		if t.Name == tag {
			jsonContent, err := json.MarshalIndent(&t, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(jsonContent))
			return nil
		}
	}
	return nil
}

func searchTag(ctx *cli.Context) error {
	repo, tag, err := getTagArgs(ctx)
	if err != nil {
		return err
	}

	tags, err := getRepoTags(globalClient, repo)
	if err != nil {
		return err
	}

	for _, t := range tags {
		if strings.Contains(t.Name, tag) {
			fmt.Println(t.Name)
		}
	}
	return nil
}

func deleteTag(ctx *cli.Context) error {
	repo, tag, err := getTagArgs(ctx)
	if err != nil {
		return err
	}

	return globalClient.DeleteTag(repo, tag)
}

func getTagArgs(ctx *cli.Context) (string , string,error) {
	if ctx.NArg() < 2 {
		return "", "", fmt.Errorf("must input repo name and tag name")
	}
	return ctx.Args().Get(0), ctx.Args().Get(1), nil
}

func getRepoTags(c *client.Client, repo string) ([]client.Tag, error) {
	tags := []client.Tag{}
	page := 1
	for {
		pageTags, err := c.GetTags(repo, page, defaultPageSize)
		if err != nil {
			return nil, fmt.Errorf("get %s tags for namespace '%s' from DockerHub error: %s", repo, c.Namespce(), err.Error())
		}
		tags = append(tags, pageTags.Tags...)

		if len(pageTags.Next) == 0 {
			break
		}
		page++
	}
	return tags, nil
}