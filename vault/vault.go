package vault

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/kildevaeld/go-filecrypt"
	"github.com/rakyll/magicmime"
)

// var (
// 	NotUniqueErr = errors.New("")
//)

type ByteCounter struct {
	written uint64
}

func (b *ByteCounter) Write(by []byte) (int, error) {
	l := len(by)
	b.written += uint64(l)
	return l, nil
}

func (b *ByteCounter) Written() uint64 {
	return b.written
}

type Item struct {
	Id        FileId
	Name      string
	Tags      []string
	Mime      string
	Size      uint64
	Encrypted bool
}

type ItemCreateOptions struct {
	Name      string
	Mime      string
	Size      uint64
	Key       *[32]byte
	Overwrite bool
	Tags      []string
}

var Key = filecrypt.Key

type FileId string

type FileStore interface {
	Create(r io.Reader, size uint64, item *Item) (FileId, error)
	Read(file FileId) (io.ReadCloser, error)
	Remove(file FileId) error
}

type MetaStore interface {
	Create(item *Item, options ItemCreateOptions) error
	Read(name string) (*Item, error)
	Find(nameOrTag string) ([]*Item, error)
	Get(id string) (*Item, error)
	Has(name string) bool
	Remove(name string) error
	List() []*Item
}

type Map map[string]interface{}

func (m Map) ToJSON() []byte {
	var b []byte
	b, _ = json.Marshal(m)
	return b
}

/*type Vault struct {
	fileStore FileStore
	metaStore MetaStore
}*/

func (v *Vault) AddFromPath(path string, o ItemCreateOptions) (*Item, error) {

	ext := filepath.Ext(path)
	mime := mime.TypeByExtension(ext)

	if mime == "" {
		if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
			return nil, err
		}

		mimetype, err := magicmime.TypeByFile(path)
		if err != nil {
			return nil, err
		}

		mime = mimetype

		magicmime.Close()
	}

	stat, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	o.Size = uint64(stat.Size())
	if o.Mime == "" {
		o.Mime = mime
	}

	return v.Add(file, o)
}

func validateOptions(options *ItemCreateOptions) bool {
	if options.Name == "" {
		return false
	} else if options.Mime == "" {
		return false
	} else if options.Size == 0 {
		return false
	}
	return true
}

func (v *Vault) AddFromMap(m Map, options ItemCreateOptions) (*Item, error) {

	b, e := json.Marshal(&m)

	if e != nil {
		return nil, e
	}

	options.Mime = "application/json"

	return v.AddBytes(b, options)
}

func (v *Vault) AddBytes(bs []byte, options ItemCreateOptions) (*Item, error) {
	b := bytes.NewReader(bs)
	options.Size = uint64(len(bs))
	return v.Add(b, options)
}

func (v *Vault) Get(id string) (*Item, error) {
	return v.metaStore.Get(id)
}

func (v *Vault) Open(item *Item) (io.ReadCloser, error) {
	/*if !v.metaStore.Has(item.Name) {
		return nil
	}*/

	return v.fileStore.Read(item.Id)
	//return o
}

func (v *Vault) Decrypt(dest io.Writer, item *Item, key *[32]byte) error {

	reader, e := v.Open(item)

	if e != nil {
		return e
	}

	return filecrypt.Decrypt(dest, reader, key)

}

func (v *Vault) Find(nameOrTag string) []*Item {
	i, _ := v.metaStore.Find(nameOrTag)
	return i
}

func (v *Vault) List() []*Item {
	return v.metaStore.List()
}
func (v *Vault) Remove(name string) error {
	if !v.metaStore.Has(name) {
		return errors.New("!exists")
	}

	item, _ := v.metaStore.Read(name)

	e := v.fileStore.Remove(item.Id)

	if e != nil {
		return e
	}

	return v.metaStore.Remove(name)
}

func NewVault(metaStore MetaStore, fileStore FileStore) *Vault {
	return &Vault{
		fileStore: fileStore,
		metaStore: metaStore,
	}
}
