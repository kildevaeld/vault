package filecrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/golang.org/x/crypto/nacl/secretbox"
)

const (
	NonceLength       = 24
	KeyLength         = 32
	HeaderLength      = 3
	SegmentOverheader = 2 + NonceLength + secretbox.Overhead
	MaxPackageSize    = ^uint16(0)
)

var (
	ErrNotUnique      = errors.New("not unique")
	ErrLengthMismatch = errors.New("length mismatch")
	PackageSize       = uint16(1024 * 32)
)

func fixedHeader(byts []byte) [HeaderLength]byte {
	if len(byts) != HeaderLength {
		panic(ErrLengthMismatch)
	}

	var out [HeaderLength]byte
	for i, b := range byts {
		out[i] = b
	}

	return out
}

func fixedKey(key []byte) [KeyLength]byte {
	if len(key) != KeyLength {
		panic(ErrLengthMismatch)
	}
	var k [KeyLength]byte
	for i, b := range key {
		k[i] = b
	}
	return k
}

func fixedNonce(nonce []byte) [NonceLength]byte {
	if len(nonce) != NonceLength {
		panic(ErrLengthMismatch)
	}

	var k [NonceLength]byte
	for i, b := range nonce {
		k[i] = b
	}
	return k
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
			err = ErrLengthMismatch
		} else if contains(used, nonce) {
			err = ErrNotUnique
		} else {
			err = nil
			break
		}

	}

	return nonce, err
}

func Key(input []byte) [32]byte {
	hash := sha256.New()
	hash.Write(input)
	k := hash.Sum(nil)
	return fixedKey(k)
}
