package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
)

var (
	//FILE = "worldcities.csv"
	FILE = "main.go"
)

func realMain() int {

	cPath, err := vault.ConfigFile()

	if err != nil {
		panic(err)
	}

	file, err := os.Open(cPath)
	var config vault.Config
	err = vault.DecodeConfig(file, &config)

	if err != nil {
		panic(err)
	}
	file.Close()

	sConf, e := server.GetServerConfig(&config)
	if e != nil {
		panic(e)
	}
	client, err := server.NewVaultClient(&server.ClientConfig{
		ServerConfig: sConf,
	})

	if err == nil {
		err = client.Ping()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to vault\n%v\n", err)
		return 1
	}

	app := cli.NewApp()
	app.Name = "vault"

	app.Version = "0.0.1"

	app.Commands = initCommands(client)

	err = app.Run(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(realMain())
}

func initCommands(client *server.Client) []cli.Command {
	return []cli.Command{
		addCommand(client),
		listCommand(client),
		mountCommand(client),
		getCommand(client),
		removeCommand(client),
	}
}

func printErrorAndExit(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}
