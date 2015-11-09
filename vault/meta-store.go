package vault

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/ryanuber/go-glob"
)

type itemMap map[FileId]*Item

func findByName(haystack itemMap, name string) (item *Item) {
	for _, v := range haystack {
		if v.Name == name {
			item = v
			break
		}
	}
	return item
}

func findByNameOrTag(haystack itemMap, nameOrTag string) []*Item {
	var out []*Item

	for _, v := range haystack {
		if glob.Glob(nameOrTag, v.Name) {
			out = append(out, v)
		}
		/*if v.Name == nameOrTag || contains(v.Tags, nameOrTag) {

			out = append(out, v)
		}*/
	}
	return out
}

func contains(haustack []string, needle string) bool {
	for _, s := range haustack {
		if s == needle {
			return true
		}
	}
	return false
}

func hasItem(haystack itemMap, name string) bool {
	return findByName(haystack, name) != nil
}

type FileSystemMetaStoreConfig struct {
	Path string
}

func (self FileSystemMetaStoreConfig) Type() string {
	return "Filesystem"
}

type FileSystemMetaStore struct {
	path  string
	files map[FileId]*Item
	lock  *sync.Mutex
}

func (f *FileSystemMetaStore) getPath(component string) string {
	return filepath.Join(f.path, component)
}

func (f *FileSystemMetaStore) Sync() (err error) {

	if !Exists(f.path) || !IsDir(f.path) {
		return errors.New("destination path does not exists!")
	}

	configFile := filepath.Join(f.path, "__store.json")

	// load
	if f.files == nil {
		if Exists(configFile) {
			b, e := ioutil.ReadFile(configFile)
			if e != nil {
				err = e
			} else {
				err = json.Unmarshal(b, &f.files)
			}
		} else {
			f.files = make(itemMap)
		}

	} else {
		b, e := json.Marshal(&f.files)
		if e != nil {
			err = e
		} else {
			var out bytes.Buffer
			json.Indent(&out, b, "", "  ")
			err = ioutil.WriteFile(configFile, out.Bytes(), 0755)
		}
	}

	return err
}

func (f *FileSystemMetaStore) Create(item *Item, o ItemCreateOptions) error {

	f.lock.Lock()
	defer f.lock.Unlock()

	oldItem := findByName(f.files, item.Name)

	if oldItem != nil && !o.Overwrite {
		return nil
	} else if oldItem != nil {
		item.Id = oldItem.Id
		f.files[oldItem.Id] = item
	} else {
		f.files[item.Id] = item
	}

	return f.Sync()
}

func (f *FileSystemMetaStore) Read(name string) (*Item, error) {

	f.lock.Lock()
	defer f.lock.Unlock()

	return findByName(f.files, name), nil
}

func (f *FileSystemMetaStore) Has(name string) bool {
	return hasItem(f.files, name)
}

func (f *FileSystemMetaStore) Remove(id string) error {

	item, e := f.Get(id)

	if e != nil {
		return e
	}

	if item == nil {
		return errors.New("not found")
	}

	delete(f.files, item.Id)

	return f.Sync()
}

func (f *FileSystemMetaStore) List() []*Item {
	var out []*Item
	for _, v := range f.files {
		out = append(out, v)
	}
	return out
}

func (f *FileSystemMetaStore) Find(nameOrTag string) ([]*Item, error) {
	return findByNameOrTag(f.files, nameOrTag), nil
}

func (f *FileSystemMetaStore) Get(id string) (*Item, error) {
	for _, r := range f.files {
		return r, nil
	}
	return nil, nil
}

func NewFileSystemMetaStore(conf FileSystemMetaStoreConfig) (*FileSystemMetaStore, error) {
	store := &FileSystemMetaStore{
		path: conf.Path,
		lock: &sync.Mutex{},
	}
	return store, store.Sync()
}
