package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
	"github.com/mitchellh/mapstructure"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {

	config, err := loadConfig()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while loading config\n%v\n", err)
		return 1
	}
	var fileStore vault.FileStore
	var metaStore vault.MetaStore

	metaStore, err = loadMetaStore(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing meta store\n%v\n", err)
		return 1
	}

	fileStore, err = loadFileStore(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing file store\n%v\n", err)
		return 1
	}

	v := vault.NewVault(metaStore, fileStore)

	s, e := loadServer(v, config)

	if e != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing server\n%v\n", e)
		return 1
	}

	s.Listen()

	defer s.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	<-ch

	return 0
}

func loadConfig() (*vault.Config, error) {

	cPath, e := vault.ConfigFile()

	if e != nil {
		return nil, e
	}

	//cPath := filepath.Join(home, "vault.toml")

	file, err := os.Open(cPath)

	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		conf := vault.DefaultConfig()

		b, e := vault.EncodeConfig(conf)

		if e != nil {
			return nil, e
		}

		ioutil.WriteFile(cPath, b, 0755)

		return conf, nil
	} else {

	}

	defer file.Close()

	var config vault.Config
	err = vault.DecodeConfig(file, &config)

	return &config, err
}

func loadMetaStore(conf *vault.Config) (vault.MetaStore, error) {
	var typ string
	t := conf.MetaStore["Type"]

	if t == nil {
		t = conf.MetaStore["type"]
	}
	if t != nil {
		typ = t.(string)
	}

	if typ == "Filesystem" || typ == "filesystem" {
		var config vault.FileSystemMetaStoreConfig
		err := mapstructure.Decode(conf.MetaStore, &config)
		if err != nil {
			return nil, err
		}
		return vault.NewFileSystemMetaStore(config)
	}

	return nil, nil

}

func loadFileStore(conf *vault.Config) (vault.FileStore, error) {
	var typ string
	t := conf.MetaStore["Type"]

	if t == nil {
		t = conf.MetaStore["type"]
	}
	if t != nil {
		typ = t.(string)
	}

	if typ == "Filesystem" || typ == "filesystem" {
		var config vault.FileSystemFileStoreConfig
		mapstructure.Decode(conf.FileStore, &config)
		return vault.NewFileSystemFileStore(config), nil
	}

	return nil, nil
}

func loadServer(v *vault.Vault, conf *vault.Config) (*server.VaultServer, error) {
	typ := ""
	t := conf.Server["Type"]

	if t == nil {
		t = conf.Server["type"]
	}
	if t != nil {
		typ = t.(string)
	}

	var sConf interface{}
	if strings.ToLower(typ) == "tcp" {
		var config server.VaultServerTCPConfig
		err := mapstructure.Decode(conf.Server, &config)
		if err != nil {
			return nil, err
		}
		sConf = config
	} else if strings.ToLower(typ) == "unix" {
		var config server.VaultServerUnixConfig
		mapstructure.Decode(conf.Server, &config)
		sConf = config
	}
	fmt.Printf("%v", sConf)
	return server.NewVaultServer(v, sConf)

}
