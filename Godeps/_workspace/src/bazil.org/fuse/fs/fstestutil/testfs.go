package fstestutil

import (
	"os"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse/fs"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/golang.org/x/net/context"
)

// SimpleFS is a trivial FS that just implements the Root method.
type SimpleFS struct {
	Node fs.Node
}

var _ = fs.FS(SimpleFS{})

func (f SimpleFS) Root() (fs.Node, error) {
	return f.Node, nil
}

// File can be embedded in a struct to make it look like a file.
type File struct{}

func (f File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = 0666
	return nil
}

// Dir can be embedded in a struct to make it look like a directory.
type Dir struct{}

func (f Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0777
	return nil
}

// ChildMap is a directory with child nodes looked up from a map.
type ChildMap map[string]fs.Node

var _ = fs.Node(&ChildMap{})
var _ = fs.NodeStringLookuper(&ChildMap{})

func (f *ChildMap) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0777
	return nil
}

func (f *ChildMap) Lookup(ctx context.Context, name string) (fs.Node, error) {
	child, ok := (*f)[name]
	if !ok {
		return nil, fuse.ENOENT
	}
	return child, nil
}
