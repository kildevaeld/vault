package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kildevaeld/vault/server"
)

func getCommand(client *server.Client) cli.Command {
	return cli.Command{
		Name: "get",
		Action: func(ctx *cli.Context) {
			getCmd(client, ctx.Args().First())
		},
		Before: func(ctx *cli.Context) error {
			if len(ctx.Args()) == 0 {
				return errors.New("usage: vault get <id>")
			}
			return nil
		},
	}
}

func getCmd(client *server.Client, id string) {

	r, e := client.Reader(id)
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	io.Copy(os.Stdout, r)
}
