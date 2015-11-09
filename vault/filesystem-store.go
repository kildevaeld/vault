package vault

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"code.google.com/p/go-uuid/uuid"
)

func Exists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

type FileSystemFileStoreConfig struct {
	Path string
}

type FileSystemFileStore struct {
	path string
	lock *sync.Mutex
}

func (f *FileSystemFileStore) Create(r io.Reader, size uint64, item *Item) (FileId, error) {

	b := uuid.New()

	path := filepath.Join(f.path, b)

	f.lock.Lock()
	defer f.lock.Unlock()

	file, err := os.Create(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	io.Copy(file, r)

	return FileId(b), err

}

func (f *FileSystemFileStore) Remove(id FileId) error {
	path := filepath.Join(f.path, string(id))
	return os.RemoveAll(path)
}

func (f *FileSystemFileStore) Read(id FileId) (io.ReadCloser, error) {

	path := filepath.Join(f.path, string(id))

	if !Exists(path) {
		return nil, errors.New("not exists")
	}

	return os.Open(path)
}

func NewFileSystemFileStore(config FileSystemFileStoreConfig) *FileSystemFileStore {
	return &FileSystemFileStore{
		path: config.Path,
		lock: &sync.Mutex{},
	}
}
