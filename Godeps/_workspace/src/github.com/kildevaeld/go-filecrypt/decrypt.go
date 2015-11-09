package filecrypt

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/golang.org/x/crypto/nacl/secretbox"
)

const MaxUint16 = ^uint16(0)

func Decrypt(dest io.Writer, src io.Reader, key *[32]byte) (err error) {

	var header [HeaderLength]byte
	hr, eh := src.Read(header[:])

	if eh != nil {
		err = eh
		return err
	}

	if hr != HeaderLength {
		return errors.New("header")
	}

	if string(header[:]) != "vau" {
		return errors.New("fileformat")
	}

	var pkgSize uint16
	buf := make([]byte, MaxUint16)
	for {

		ep := binary.Read(src, binary.LittleEndian, &pkgSize)

		if ep != nil {
			err = ep
			break
		}

		segSize := int(pkgSize) + NonceLength + secretbox.Overhead

		nr, er := src.Read(buf[0:segSize])

		if er != nil {
			er = err
			break
		}

		if nr != segSize {
			err = errors.New("short read")
			break
		}

		msg, em := DecryptMessage(buf[0:segSize], key)

		if em != nil {
			err = em
			break
		}

		nw, ew := dest.Write(msg)

		if ew != nil {
			err = ew
			break
		}

		if nw != int(pkgSize) {
			err = io.ErrShortWrite
			break
		}

	}
	if err == io.EOF {
		err = nil
	}

	return err
}

func DecryptMessage(src []byte, key *[32]byte) ([]byte, error) {

	nonce := fixedNonce(src[0:NonceLength])

	var ok bool
	var buf []byte

	if buf, ok = secretbox.Open(buf, src[NonceLength:], &nonce, key); ok {
		return buf, nil
	}

	return nil, errors.New("decrypt")
}
