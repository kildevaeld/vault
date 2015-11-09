package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	mC, _ := config.GetMetaConfig()
	log.Printf("using metastore: %s\n", mC.Type())

	fileStore, err = loadFileStore(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing file store\n%v\n", err)
		return 1
	}

	fC, _ := config.GetFileStoreConfig()
	log.Printf("using filestore: %s\n", fC.Type())

	v := vault.NewVault(metaStore, fileStore)

	s, e := loadServer(v, config)

	if e != nil {
		fmt.Fprintf(os.Stderr, "Error while initializing server\n%v\n", e)
		return 1
	}

	s.Listen()
	//err = wait(s)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	return 0
}

func wait(s *server.VaultServer) error {
	done := make(chan error, 1)
	defer close(done)
	ch := make(chan os.Signal, 1)
	defer close(ch)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		err := s.Listen()
		if done != nil {
			done <- err
		}

	}()

	select {
	case <-ch:
		s.Close()
	case err := <-done:
		return err

	}
	return nil
	/*
		for {
			select {
			case <-ch:
				s.Close()
			case err := <-done:
				return err
			default:
				// move along
			}
		}*/
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

	config, e := conf.GetMetaConfig()

	if e != nil {
		return nil, e
	}

	var metaStore vault.MetaStore
	if fs, ok := config.(vault.FileSystemMetaStoreConfig); ok {

		metaStore, e = vault.NewFileSystemMetaStore(fs)
	}

	return metaStore, e
}

func loadFileStore(conf *vault.Config) (vault.FileStore, error) {
	config, e := conf.GetFileStoreConfig()

	if e != nil {
		return nil, e
	}

	if fs, ok := config.(vault.FileSystemFileStoreConfig); ok {
		return vault.NewFileSystemFileStore(fs), nil
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

	return server.NewVaultServer(v, sConf)

}
