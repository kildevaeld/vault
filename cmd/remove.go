package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/kildevaeld/vault/server"
)

func removeCommand(client *server.Client) cli.Command {
	return cli.Command{
		Name:    "remove",
		Aliases: []string{"r", "rm"},
		Action: func(ctx *cli.Context) {
			removeCmd(client, ctx.Args().First())
		},
		Before: func(ctx *cli.Context) error {
			if len(ctx.Args()) == 0 {
				return errors.New("usage: vault rm <id>")
			}
			return nil
		},
	}
}

func removeCmd(client *server.Client, id string) {

	err := client.Remove(id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

}
