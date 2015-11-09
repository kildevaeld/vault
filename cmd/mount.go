package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/prompt"
	"github.com/kildevaeld/vault/mount"
	"github.com/kildevaeld/vault/server"
)

func mountCommand(client *server.Client) cli.Command {

	return cli.Command{
		Name: "mount",
		Action: func(ctx *cli.Context) {
			mountpoint := ctx.Args().First()
			mountCmd(client, mountpoint)
		},
		Before: func(ctx *cli.Context) error {
			if len(ctx.Args()) == 0 {
				return errors.New("usage: vault mount <mountpoint>")
			}
			return nil
		},
	}

}

func mountCmd(client *server.Client, mountpoint string) {

	if !filepath.IsAbs(mountpoint) {

		abs, err := filepath.Abs(mountpoint)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		mountpoint = abs
	}

	if _, err := os.Lstat(mountpoint); err != nil {
		fmt.Fprintf(os.Stderr, "mountpoint does not exists: %s", mountpoint)
		os.Exit(1)
	}

	m := mount.NewMount(client, mountpoint)

	done := make(chan bool, 1)
	defer close(done)
	go func() {
		err := m.Serve()

		if err != nil {
			fmt.Printf("%v\n", err)
		}
		done <- true
	}()

	cleanup := func() {
		err := prompt.NewProcess("Umounting "+mountpoint+" ...", func() error {
			return m.Close()
		})
		ret := 0
		if err != nil {
			fmt.Printf("%v\n", err)
			ret = 1
		}
		os.Exit(ret)
	}

	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		select {
		case <-done:
			cleanup()
		case <-ch:
			cleanup()
		default:
			// just move along
		}

	}

}
