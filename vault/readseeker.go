package vault

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync"
)

type ReadWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Closer
	io.Seeker
}

type FileBuffer struct {
	Buffer *bytes.Buffer
	Index  int64
}

func NewFileBuffer(b []byte) *FileBuffer {
	if b == nil {
		b = make([]byte, 1)
	}
	return &FileBuffer{bytes.NewBuffer(b), 0}
}

func (fbuffer *FileBuffer) Bytes() []byte {
	return fbuffer.Buffer.Bytes()
}

func (fBuffer *FileBuffer) Size() int64 {
	return int64(len(fBuffer.Bytes()))
}

func (fbuffer *FileBuffer) Read(p []byte) (int, error) {

	if fbuffer.Index == fbuffer.Size() {
		return 0, io.EOF
	}

	n, err := bytes.NewBuffer(fbuffer.Buffer.Bytes()[fbuffer.Index:]).Read(p)

	if err == nil {
		if fbuffer.Index+int64(len(p)) < int64(fbuffer.Buffer.Len()) {
			fbuffer.Index += int64(len(p))
		} else {
			fbuffer.Index = int64(fbuffer.Buffer.Len())
		}
	}

	return n, err
}

func (fbuffer *FileBuffer) Write(p []byte) (int, error) {

	if fbuffer.Index == fbuffer.Size() {
		n, err := fbuffer.Buffer.Write(p)

		if err == nil {
			fbuffer.Index = int64(fbuffer.Buffer.Len())
		}
		return n, err
	} else {
		//lng := fbuffer.Size()
		b := fbuffer.Buffer.Bytes()
		b = append(b[:fbuffer.Index], append(p, b[fbuffer.Index:]...)...)
		fbuffer.Buffer = bytes.NewBuffer(b)

		newOffset := fbuffer.Index + int64(len(p))

		fbuffer.Index = newOffset

		return len(p), nil
	}

}

func (fbuffer *FileBuffer) Seek(offset int64, whence int) (int64, error) {
	var err error
	var Index int64 = 0
	lng := fbuffer.Size()
	switch whence {
	case 0:
		if offset >= lng || offset < 0 {
			err = errors.New("Invalid Offset.")
		} else {
			fbuffer.Index = offset
			Index = offset
		}
	case 2:
		if offset > lng || lng-offset < 0 {
			err = errors.New("Invalid Offset ")

		} else {
			fbuffer.Index = lng - offset
			Index = fbuffer.Index
		}
	case 1:
		if offset+fbuffer.Index > lng {
			err = errors.New("Invalid Offset")
		} else {
			fbuffer.Index = offset + fbuffer.Index
			Index = fbuffer.Index
		}
	default:
		err = errors.New("Unsupported Seek Method.")
	}

	return Index, err
}

func (f *FileBuffer) Close() error {
	f.Buffer.Reset()
	return nil
}

type ReadSeeker struct {
	reader   io.Reader
	buffer   ReadWriteSeekCloser
	position int64
	read     int64
	haseof   bool
	mu       sync.Mutex
}

func (self *ReadSeeker) Size() int64 {
	self.mu.Lock()
	defer self.mu.Unlock()

	pos := self.position

	n, _ := self.Seek(0, 2)

	self.Seek(pos, 0)

	return n

}

func (self *ReadSeeker) Seek(offset int64, whence int) (int64, error) {

	self.mu.Lock()
	defer self.mu.Unlock()

	if whence == 0 {
		if offset > self.read {
			b := make([]byte, offset-self.read)
			read, err := self.reader.Read(b)

			if err != nil {
				return -1, err
			}

			self.buffer.Write(b)

			self.read += int64(read)
			self.position = int64(read)
			return self.position, nil
		} else {
			self.position = offset
			return self.buffer.Seek(offset, whence)
		}

	} else if whence == 1 {
		//newOffset := offset + self.position

	} else if whence == 2 {
		if self.haseof {
			p, e := self.buffer.Seek(offset, whence)
			if e != nil {
				return 0, e
			}
			self.position = p
			return p, e
		} else {

			_, e := self.buffer.Seek(0, 2)
			if e != nil {
				return 0, e
			}
			p, e := io.Copy(self.buffer, self.reader)

			if e != nil {
				return p, e
			}
			self.read += p

			p, e = self.buffer.Seek(offset, 2)

			if e != nil {
				return p, e
			}

			self.position = p
			self.haseof = true

			return p, e
		}
	}

	return 0, nil
}

func (self *ReadSeeker) Read(b []byte) (int, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.position < self.read {
		l := int64(len(b))

		// read from buffer
		if l <= self.read-self.position {

			read, err := self.buffer.Read(b)

			if err != nil {
				if err == io.EOF {
					self.haseof = true
				} else {
					return 0, err
				}

			}

			self.position += int64(read)
			fmt.Printf("read buffer: %#v\n", self.buffer)
			return read, nil
		} else {
			diff := l - (self.read - self.position)
			fmt.Printf("diff %v", diff)

		}

	} else {
		//l := int64(len(b))
		read, err := self.reader.Read(b)

		if err != nil {
			if err == io.EOF {
				self.haseof = true
			} else {
				return read, err
			}

		}
		self.buffer.Write(b)
		self.position += int64(read)
		self.read += int64(read)
		//fmt.Printf("READ %v", self)
		return read, err
	}

	/*l := len(b)
	diff := read - position
	if diff > 0 {

	}*/

	return 0, nil
}

func (self *ReadSeeker) Close() error {
	return self.buffer.Close()
}

func (self *ReadSeeker) Write(b []byte) (int, error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	n, e := self.buffer.Write(b)
	if e != nil {
		lng := int64(n)
		newOffset := lng + self.position
		if newOffset > self.read {
			self.read = newOffset
		}
		self.position = newOffset
	}

	return n, e
}

func NewReadSeeker(r io.Reader, buffer ReadWriteSeekCloser) *ReadSeeker {
	return &ReadSeeker{r, buffer, 0, 0, false, sync.Mutex{}}
}
