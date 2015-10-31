package vault

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

type Header struct {
	Version  byte
	Packages int
}

func (h *Header) Serialize() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, h.Version)
	binary.Write(buf, binary.LittleEndian, h.Packages)
	return buf.Bytes()
}

func (h *Header) Deserialize(b []byte) error {

	return nil
}

var (
	NotUniqueErr      = errors.New("Uninque")
	LengthMisMatchErr = errors.New("Length")
)

const (
	KeyLength   = 32
	NonceLength = 24
	PackageSize = 1024 * 32
	SegmentSize = PackageSize + NonceLength + secretbox.Overhead
)

func fixedKey(key []byte) [KeyLength]byte {
	if len(key) > KeyLength || len(key) < KeyLength {
		panic(LengthMisMatchErr)
	}
	var k [KeyLength]byte
	for i, b := range key {
		k[i] = b
	}
	return k
}

func fixedNonce(nonce []byte) [NonceLength]byte {
	if len(nonce) > NonceLength || len(nonce) < NonceLength {
		panic(LengthMisMatchErr)
	}

	var k [NonceLength]byte
	for i, b := range nonce {
		k[i] = b
	}
	return k
}

func Key(input []byte) [32]byte {
	hash := sha256.New()
	hash.Write(input)
	k := hash.Sum(nil)
	return fixedKey(k)
}

func serialize(position int32, nonce [24]byte, box *[]byte) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, SegmentSize))

	//binary.Write(buf, binary.LittleEndian, &position)
	buf.Write(nonce[:])
	buf.Write(*box)

	return buf.Bytes()
}

func Encrypt(dest io.Writer, src io.Reader, length int64, key [32]byte) (err error) {
	bufSize := int64(PackageSize)
	buf := make([]byte, bufSize)

	pkgs := length / bufSize
	last := length % bufSize

	var nonces [][]byte

	//fmt.Printf("Packages %d, last: %d, diff :%d\n", pkgs, last, bufSize-last)

	np := PackageSize + secretbox.Overhead + 24
	var written int
	for i := int64(0); i < pkgs; i++ {

		nr, er := src.Read(buf)

		if er != nil {
			err = er
			break
		}

		if nr < int(bufSize) {
			err = LengthMisMatchErr
			break
		}

		msg, me := EncryptMessage(buf, &key, int32(i), &nonces)

		if me != nil {
			return me
		}

		nw, ew := dest.Write(msg)

		if nw > 0 {
			written += nw
		}
		if ew != nil {
			err = ew
			break
		}

		if nw != np {

			err = LengthMisMatchErr
			break
		}

	}

	if err != nil {
		return err
	}

	if last > 0 {
		nr, er := src.Read(buf)

		if er != nil {
			return er
		}
		if int64(nr) != last {
			return LengthMisMatchErr
		}

		for i := last; i < bufSize; i++ {
			buf[i] = 0x0
		}

		msg, me := EncryptMessage(buf, &key, int32(pkgs+1), &nonces)

		if me != nil {
			return me
		}

		nw, ew := dest.Write(msg)

		if nw > 0 {
			written += nw
		}
		if ew != nil {
			err = ew
		}

		if nw != np {
			err = LengthMisMatchErr
		}

	}

	return err

}

func EncryptMessage(src []byte, key *[32]byte, position int32, nonces *[][]byte) ([]byte, error) {

	n, en := GetNunce(NonceLength, 10, nonces)

	nonce := fixedNonce(n)

	if en != nil {
		return nil, en
	}

	var dest []byte
	dest = secretbox.Seal(dest, src, &nonce, key)

	return serialize(position, nonce, &dest), nil

}

func DecryptMessage(src []byte, key *[32]byte, position int32) ([]byte, error) {
	if len(src) != SegmentSize {
		return nil, LengthMisMatchErr
	}

	nonce := fixedNonce(src[0:NonceLength])

	var ok bool
	buf := make([]byte, PackageSize)

	if buf, ok = secretbox.Open(buf, src[NonceLength:], &nonce, key); ok {
		return buf, nil
	}

	return nil, errors.New("decrypt")
}

func Decrypt(dest io.Writer, src io.Reader, length int64, key [32]byte) (err error) {

	pkgsSize := int64(SegmentSize)
	pkgs := length / pkgsSize

	buf := make([]byte, pkgsSize)

	if length%pkgsSize > 0 {
		return LengthMisMatchErr
	}

	for i := int64(0); i < pkgs; i++ {

		nr, er := src.Read(buf)

		if er != nil {
			err = er
			break
		}
		if int64(nr) != pkgsSize {
			err = LengthMisMatchErr
			break
		}

		msg, em := DecryptMessage(buf, &key, int32(i))
		if em != nil {
			err = em
			break
		}

		nw, ew := dest.Write(msg)

		if ew != nil {
			err = ew
			break
		}
		if nw != len(msg) {
			err = LengthMisMatchErr
			break
		}

	}

	return err
}

func contains(haystack *[][]byte, needle []byte) bool {
	for _, h := range *haystack {
		if bytes.Compare(h, needle) == 0 {
			return true
		}
	}
	return false
}

func GetNunce(length, retries int, used *[][]byte) ([]byte, error) {
	nonce := make([]byte, length)

	var err error
	for i := 0; i < retries; i++ {
		l, e := rand.Read(nonce[:])

		if e != nil {
			err = e
			break
		}

		if l != length {
			err = LengthMisMatchErr
		} else if contains(used, nonce) {
			err = NotUniqueErr
		} else {
			err = nil
			break
		}

	}

	return nonce, err
}
