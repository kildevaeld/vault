package mount

import (
	"io"
	"sync"
	"time"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/bazil.org/fuse/fs"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/kildevaeld/vault/server"
	"github.com/kildevaeld/vault/vault"
)

type File struct {
	client *server.Client
	item   *vault.Item
	buffer io.ReadWriteCloser
	mu     sync.Mutex
	seek   *vault.ReadSeeker
}

var _ fs.Node = (*File)(nil)

func (self *File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = 0755
	a.Size = self.item.Size
	t := time.Now()
	a.Crtime = t
	a.Ctime = t
	a.Mtime = t

	return nil
}

var _ = fs.NodeOpener(&File{})

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.seek == nil {
		r, err := f.client.Reader(string(f.item.Id))
		if err != nil {
			return nil, err
		}
		buffer := vault.NewFileBuffer(nil)
		seek := vault.NewReadSeeker(r, buffer)
		f.seek = seek
	}

	// individual entries inside a zip file are not seekable
	//resp.Flags |= fuse.OpenNonSeekable
	return &FileHandle{r: f.seek}, nil
}

type FileHandle struct {
	r  *vault.ReadSeeker
	mu sync.Mutex
}

var _ fs.Handle = (*FileHandle)(nil)

var _ fs.HandleReleaser = (*FileHandle)(nil)

func (fh *FileHandle) Release(ctx context.Context, req *fuse.ReleaseRequest) error {

	return fh.r.Close()
}

var _ = fs.HandleReader(&FileHandle{})

func (fh *FileHandle) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	// We don't actually enforce Offset to match where previous read
	// ended. Maybe we should, but that would mean'd we need to track
	// it. The kernel *should* do it for us, based on the
	// fuse.OpenNonSeekable flag.
	//
	// One exception to the above is if we fail to fully populate a
	// page cache page; a read into page cache is always page aligned.
	// Make sure we never serve a partial read, to avoid that.
	_, e := fh.r.Seek(req.Offset, 0)
	if e != nil {
		return e
	}
	//buf := bytes.NewBuffer(nil)

	buf := make([]byte, req.Size)
	n, err := fh.r.Read(buf) //io.ReadFull(fh.r, buf)
	if err == io.ErrUnexpectedEOF || err == io.EOF {
		err = nil
	}
	resp.Data = buf[:n]
	return err

}

var _ = fs.HandleWriter(&FileHandle{})

const maxInt = int(^uint(0) >> 1)

func (f *FileHandle) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// expand the buffer if necessary
	/*newLen := req.Offset + int64(len(req.Data))
	if newLen > int64(maxInt) {
		return fuse.Errno(syscall.EFBIG)
	}
	if newLen := int(newLen); newLen > len(f.data) {
		f.data = append(f.data, make([]byte, newLen-len(f.data))...)
	}

	n := copy(f.data[req.Offset:], req.Data)*/
	n, e := f.r.Write(req.Data)
	if e != nil {
		return e
	}
	resp.Size = n
	return nil
}

/*
var _ = fs.HandleFlusher(&File{})

func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
  f.mu.Lock()
  defer f.mu.Unlock()

  if f.writers == 0 {
    // Read-only handles also get flushes. Make sure we don't
    // overwrite valid file contents with a nil buffer.
    return nil
  }

  err := f.dir.fs.db.Update(func(tx *bolt.Tx) error {
    b := f.dir.bucket(tx)
    if b == nil {
      return fuse.ESTALE
    }
    return b.Put(f.name, f.data)
  })
  if err != nil {
    return err
  }
  return nil
}

var _ = fs.NodeSetattrer(&File{})

func (f *File) Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error {
  f.mu.Lock()
  defer f.mu.Unlock()

  if req.Valid.Size() {
    if req.Size > uint64(maxInt) {
      return fuse.Errno(syscall.EFBIG)
    }
    newLen := int(req.Size)
    switch {
    case newLen > len(f.data):
      f.data = append(f.data, make([]byte, newLen-len(f.data))...)
    case newLen < len(f.data):
      f.data = f.data[:newLen]
    }
  }
  return nil
}
*/
