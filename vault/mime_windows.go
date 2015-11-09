// +build windows
package vault

import "bitbucket.org/taruti/mimemagic"

func detectContentType(sample []byte) (string, error) {
	return mimemagic.Match("", sample), nil
}
