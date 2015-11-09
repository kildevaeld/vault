package filecrypt

import (
	"encoding/binary"
	"io"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/golang.org/x/crypto/nacl/secretbox"
)

func writeHeader(w io.Writer) error {

	nw, err := w.Write([]byte("vau"))

	if nw != HeaderLength && err == nil {
		err = io.ErrShortWrite
	}

	return err

}

func Encrypt(dest io.Writer, src io.Reader, key *[32]byte) (written uint64, err error) {
	buf := make([]byte, PackageSize)

	// write the header
	if e := writeHeader(dest); e != nil {
		return 0, e
	}

	var nonces [][]byte

	for {

		nr, er := src.Read(buf)

		if nr > 0 {

			msg, me := EncryptMessage(buf[0:nr], key, &nonces)

			if me != nil {
				err = me
			}

			// Write package length
			if e := binary.Write(dest, binary.LittleEndian, uint16(nr)); e != nil {
				err = e
				break
			}
			written += 2 // package length int16

			nw, ew := dest.Write(msg)

			written += uint64(nw)

			if ew != nil {
				err = ew
				break
			}

			if nw != len(msg) {
				err = io.ErrShortWrite
				break
			}

		}

		if er == io.EOF {
			break
		}

		if er != nil {
			err = er
			break
		}

	}

	return written, err

}

func EncryptMessage(src []byte, key *[32]byte, nonces *[][]byte) ([]byte, error) {
	if nonces == nil {
		nonces = new([][]byte)
	}

	n, en := GetNunce(NonceLength, 10, nonces)

	nonce := fixedNonce(n)

	if en != nil {
		return nil, en
	}

	var dest []byte
	dest = secretbox.Seal(dest, src, &nonce, key)

	dest = append(n[:], dest...)

	return dest, nil

}
