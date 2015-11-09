// +build darwin freebsd linux netbsd openbsd
package vault

import "github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/rakyll/magicmime"

func detectContentType(sample []byte) (string, error) {
	if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
		return "", err
	}
	defer magicmime.Close()

	return magicmime.TypeByBuffer(sample)

}
