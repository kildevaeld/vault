package vault

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type Item struct {
	Id   int
	Name string
	Tags []string
	Mime string
	Size int64
}

type ItemCreate struct {
	Name   string
	Mime   string
	Size   int64
	Reader io.Reader
}

type FileStore interface {
	Create(item *ItemCreate) error
	Read(name string) (io.ReadCloser, error)
}

type MetaStore interface {
}

type Vault struct {
	fileStore FileStore
	metaStore MetaStore
}

func (v *Vault) AddFromPath(name, path string) (*Item, error) {

	ext := filepath.Ext(path)
	mime := mime.TypeByExtension(ext)

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	return v.Add(name, file, mime, 0)
}

func (v *Vault) AddFromMap(name string, m map[string]interface{}) (*Item, error) {

	b, e := json.Marshal(&m)

	if e != nil {
		return nil, e
	}

	return v.AddBytes(name, b, "application/json")
}

func (v *Vault) AddBytes(name string, bs []byte, mime string) (*Item, error) {
	b := bytes.NewReader(bs)
	return v.Add(name, b, mime, int64(len(bs)))
}

func (v *Vault) Add(name string, reader io.Reader, mime string, size int64) (*Item, error) {

	item := ItemCreate{
		Name:   name,
		Mime:   mime,
		Size:   size,
		Reader: reader,
	}

	v.fileStore.Create(&item)

	i := Item{
		Name: name,
		Mime: mime,
		Size: size,
	}

	return &i, nil
}

func (v *Vault) Get(name string) *Item {
	return nil
}

func (v *Vault) Open(item *Item) io.Reader {
	o, _ := v.fileStore.Read(item.Name)
	return o
}

func NewVault() *Vault {
	return &Vault{
		fileStore: NewFileSystemStore("./test"),
	}
}
