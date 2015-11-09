package mount

import (
	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse/fs"
	"github.com/kildevaeld/vault/server"
)

type Mount struct {
	conn   *fuse.Conn
	client *server.Client
	mount  string
}

func (self *Mount) Serve() error {

	c, err := fuse.Mount(self.mount, fuse.FSName("vault-mount"), fuse.VolumeName("vault"))
	if err != nil {
		return err
	}

	self.conn = c

	filesys := &FS{
		client: self.client,
	}
	if err := fs.Serve(c, filesys); err != nil {
		return err
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		return err
	}
	return nil
}

func (self *Mount) Close() error {
	err := self.conn.Close()

	if err != nil {
		return err
	}

	return fuse.Unmount(self.mount)

}

func NewMount(client *server.Client, mountpoint string) *Mount {
	return &Mount{
		mount:  mountpoint,
		client: client,
	}
}
