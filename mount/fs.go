package mount

import (
	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse/fs"
	"github.com/kildevaeld/vault/server"
)

type FS struct {
	client *server.Client
}

var _ fs.FS = (*FS)(nil)

func (f *FS) Root() (fs.Node, error) {
	n := &Dir{
		client: f.client,
	}
	return n, nil
}
