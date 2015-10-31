package vault

import (
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

type FileSystemStore struct {
	path string
	lock *sync.Mutex
}

func (f *FileSystemStore) Create(item *ItemCreate) error {

	path := filepath.Join(f.path, item.Name)
	f.lock.Lock()
	defer f.lock.Unlock()

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()
	var key [32]byte
	rand.Reader.Read(key[:])
	//_, err = Encrypt(file, item.Reader, key)
	//fmt.Printf("KEY %v", string(key[:]))
	return err
}

func (f *FileSystemStore) Read(name string) (io.ReadCloser, error) {

	path := filepath.Join(f.path, name)

	if !Exists(path) {
		return nil, errors.New("not exists")
	}

	return os.Open(path)
}

func NewFileSystemStore(path string) *FileSystemStore {
	return &FileSystemStore{
		path: path,
		lock: &sync.Mutex{},
	}
}
