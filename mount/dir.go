package mount

import (
	"archive/zip"
	"os"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse/fs"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/golang.org/x/net/context" // nil for the root directory, which has no entry in the zip
	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
)

type Dir struct {
	client *server.Client

	item *vault.Item
}

var _ fs.Node = (*Dir)(nil)

func zipAttr(f *zip.File, a *fuse.Attr) {
	a.Size = f.UncompressedSize64
	a.Mode = f.Mode()
	a.Mtime = f.ModTime()
	a.Ctime = f.ModTime()
	a.Crtime = f.ModTime()
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	if d.item == nil {
		// root directory
		a.Mode = os.ModeDir | 0755
		return nil
	}
	//zipAttr(d.file, a)
	return nil
}

var _ = fs.NodeRequestLookuper(&Dir{})

func (self *Dir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	path := req.Name
	if self.item != nil {
		path = self.item.Name + path
	}

	items, err := self.client.List()

	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Name == path {
			child := &File{
				client: self.client,
				item:   item,
			}

			return child, nil
		}
	}

	return nil, fuse.ENOENT
}

var _ = fs.HandleReadDirAller(&Dir{})

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	/*prefix := ""
	if d.item != nil {
		prefix = d.item.Name
	}*/

	var res []fuse.Dirent

	items, err := d.client.List()

	if err != nil {
		return nil, err
	}

	for _, item := range items {

		if item.Encrypted {
			continue
		}

		de := fuse.Dirent{
			Name: item.Name,
		}
		res = append(res, de)
	}

	/*for _, f := range d.archive.File {
		if !strings.HasPrefix(f.Name, prefix) {
			continue
		}
		name := f.Name[len(prefix):]
		if name == "" {
			// the dir itself, not a child
			continue
		}
		if strings.ContainsRune(name[:len(name)-1], '/') {
			// contains slash in the middle -> is in a deeper subdir
			continue
		}
		var de fuse.Dirent
		if name[len(name)-1] == '/' {
			// directory
			name = name[:len(name)-1]
			de.Type = fuse.DT_Dir
		}
		de.Name = name
		res = append(res, de)
	}*/
	return res, nil
}
