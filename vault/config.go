package vault

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	Server    map[string]interface{}
	MetaStore map[string]interface{}
	FileStore map[string]interface{}
}

func DecodeConfig(r io.Reader, config *Config) error {

	_, e := toml.DecodeReader(r, config)

	if e != nil {
		return e
	}

	return nil
}

func EncodeConfig(config *Config) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)

	err := encoder.Encode(config)

	return buf.Bytes(), err

}

func DefaultConfig() *Config {
	conf := Config{}

	cPath, e := ConfigDir()

	if e != nil {
		panic(e)
	}

	dir := filepath.Join(cPath, "file-store")

	_, err := os.Stat(dir)

	if err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil
		}
	}

	conf.MetaStore = Map{
		"Type": "Filesystem",
		"Path": dir,
	}

	conf.FileStore = Map{
		"Type": "Filesystem",
		"Path": dir,
	}

	conf.Server = Map{
		"Type": "unix",
		"Path": filepath.Join(cPath, "vault.socket"),
	}

	return &conf

}

func (v *Config) GetMetaConfig() (MetaStoreConfig, error) {
	var typ string
	t := v.MetaStore["Type"]

	if t == nil {
		t = v.MetaStore["type"]
	}
	if t != nil {
		typ = t.(string)
	}

	if typ == "Filesystem" || typ == "filesystem" {
		var config FileSystemMetaStoreConfig
		err := mapstructure.Decode(v.MetaStore, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return nil, errors.New("meta")

}

func (v *Config) GetFileStoreConfig() (FileStoreConfig, error) {
	var typ string
	t := v.FileStore["Type"]

	if t == nil {
		t = v.FileStore["type"]
	}
	if t != nil {
		typ = t.(string)
	}

	if typ == "Filesystem" || typ == "filesystem" {
		var config FileSystemFileStoreConfig
		mapstructure.Decode(v.FileStore, &config)
		return config, nil
	}

	return nil, errors.New("filestore")

}

func ConfigDir() (string, error) {
	home, e := homedir.Dir()

	if e != nil {
		return "", e
	}

	cPath := filepath.Join(home, ".vault")

	if !FileExists(cPath) {
		err := os.MkdirAll(cPath, 0744)

		if err != nil {
			return "", err
		}
	}

	return cPath, nil
}

func ConfigFile() (string, error) {

	cPath, err := ConfigDir()

	if err != nil {
		return "", err
	}

	cPath = filepath.Join(cPath, "vault.toml")

	if !FileExists(cPath) {
		c := DefaultConfig()
		b, e := EncodeConfig(c)
		if e != nil {
			return "", e
		}
		ioutil.WriteFile(cPath, b, 0755)
	}

	return cPath, nil

}
