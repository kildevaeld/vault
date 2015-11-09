package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
)

type addConfig struct {
	name string
	json bool
	key  string
}

func addCommand(client *server.Client) cli.Command {
	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Action: func(ctx *cli.Context) {
			file := ctx.Args().First()
			conf := addConfig{}
			conf.name = ctx.String("name")
			conf.json = ctx.Bool("json")
			conf.key = ctx.String("key")

			AddFile(file, conf, client)

		},
		Before: func(ctx *cli.Context) error {
			if len(ctx.Args()) == 0 {
				return errors.New("usage: vault add <file>")
			}
			return nil
		},
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name: "json, j",
			},
			cli.StringFlag{
				Name: "key, k",
			},
			cli.StringFlag{
				Name: "name, n",
			},
		},
	}
}

func AddFile(path string, config addConfig, client *server.Client) {

	item, e := add(path, config, client)

	if e != nil {
		fmt.Fprintf(os.Stderr, "%s\n", e.Error())
		os.Exit(1)
	}

	fmt.Printf("Added file: %v\n", item.Id)
}

func add(path string, config addConfig, client *server.Client) (*vault.Item, error) {

	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)

		if err != nil {
			return nil, err
		}

		path = abs

	}

	origiNalPath := path

	stat, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if stat.IsDir() {

		tmpDir, err := ioutil.TempDir("", "vault")

		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(tmpDir)
		tmpDir = filepath.Join(tmpDir, "vault_zip.zip")

		err = vault.Pack(path, tmpDir)

		if err != nil {
			return nil, err
		}

		path = tmpDir

	}
	if config.name == "" {

		basename := filepath.Base(origiNalPath)
		config.name = basename
	}

	c := vault.ItemCreateOptions{
		Name: config.name,
		Size: uint64(stat.Size()),
	}

	if config.key != "" {
		key := vault.Key([]byte(config.key))
		c.Key = &key
	}

	if c.Mime == "" {
		ext := filepath.Ext(path)
		c.Mime = mime.TypeByExtension(ext)
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	return client.Create(file, c)
}
