package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/dustin/go-humanize"
	"github.com/kildevaeld/prompt"
	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
)

func listCommand(client *server.Client) cli.Command {
	return cli.Command{
		Name:    "list",
		Aliases: []string{"ls", "l"},
		Action: func(ctx *cli.Context) {
			ListFiles(client)
		},
	}
}

func ListFiles(client *server.Client) {

	var items []*vault.Item
	var err error
	err = prompt.NewProcess("Retriving file information ...", func() error {
		items, err = client.List()

		return err
	})

	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	fmt.Printf("total %d\n", len(items))

	writer := tabwriter.NewWriter(os.Stdout, 1, 10, 0, '\t', 0)
	for _, i := range items {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", humanize.Bytes(i.Size), i.Name, i.Mime, i.Id)

	}
	writer.Flush()

}
