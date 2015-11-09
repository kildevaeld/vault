package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/dustin/go-humanize"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/googollee/gocmd"
	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
)

func replCommand(client *server.Client) cli.Command {
	return cli.Command{
		Name: "repl",
		Action: func(ctx *cli.Context) {
			replCmd(client)
		},
	}
}

type replServer struct {
	client *server.Client
}

func (self *replServer) Ls(server gocmd.CmdServer, args []string) {
	var items []*vault.Item
	var err error

	if len(args) > 0 {
		items, err = self.client.Find(args[0])
	} else {
		items, err = self.client.List()
	}

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("total %d\n", len(items))

	writer := tabwriter.NewWriter(os.Stdout, 1, 8, 1, '\t', 0)
	for _, i := range items {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", humanize.Bytes(i.Size), i.Name, i.Mime, i.Id)

	}
	writer.Flush()

}

func (self *replServer) Exit(server gocmd.CmdServer, args []string) {
	server.Exit(0)
}

func (self *replServer) Help(server gocmd.CmdServer, args []string) {
	fmt.Println(`
usage:
  help  show this message
  ls    list content
    `)
}

func replCmd(client *server.Client) {

	s := gocmd.NewServer("vault> ")
	s.Register(&replServer{client})
	s.Serve()

}
