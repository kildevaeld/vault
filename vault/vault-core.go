package vault

import (
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/kildevaeld/go-filecrypt"
)

type Vault struct {
	fileStore FileStore
	metaStore MetaStore
}

func (v *Vault) Add(reader io.Reader, o ItemCreateOptions) (*Item, error) {

	if v.metaStore.Has(o.Name) && !o.Overwrite {
		return nil, errors.New("already exists")
	}

	if !validateOptions(&o) {
		return nil, errors.New("options error")
	}

	i := Item{
		Name: o.Name,
		Mime: o.Mime,
		Size: o.Size,
		Tags: o.Tags,
	}

	size := o.Size

	if o.Key != nil {

		file, err := ioutil.TempFile(os.TempDir(), "prefix")
		defer os.Remove(file.Name())

		byteCounter := &ByteCounter{0}
		tee := io.TeeReader(reader, byteCounter)

		size, err = filecrypt.Encrypt(file, tee, o.Key)

		i.Size = byteCounter.Written()

		if err != nil {
			return nil, err
		}

		_, err = file.Seek(0, 0)

		if err != nil {
			return nil, err
		}

		reader = file

		i.Encrypted = true

	}

	id, e := v.fileStore.Create(reader, size, &i)

	if e != nil {
		return nil, e
	}

	i.Id = id

	err := v.metaStore.Create(&i, o)

	if err != nil {
		return nil, v.fileStore.Remove(i.Id)
	}

	return &i, err
}
