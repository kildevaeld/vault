package vault

import (
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

type FileStoreConfig interface {
	Type() string
}

type MetaStore interface {
	Create(item *Item, options ItemCreateOptions) error
	Read(name string) (*Item, error)
	Find(nameOrTag string) ([]*Item, error)
	Get(id string) (*Item, error)
	Has(name string) bool
	Remove(id string) error
	List() []*Item
}

type MetaStoreConfig interface {
	Type() string
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

func (v *Vault) Get(id string) (*Item, error) {
	return v.metaStore.Get(id)
}

func (v *Vault) Open(item *Item) (io.ReadCloser, error) {

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
func (v *Vault) Remove(id string) error {

	item, e := v.metaStore.Get(id)

	if e != nil {
		return e
	}
	if item == nil {
		return errors.New("!exists")
	}

	e = v.fileStore.Remove(item.Id)

	if e != nil {
		return e
	}

	return v.metaStore.Remove(id)
}

func NewVault(metaStore MetaStore, fileStore FileStore) *Vault {
	return &Vault{
		fileStore: fileStore,
		metaStore: metaStore,
	}
}
