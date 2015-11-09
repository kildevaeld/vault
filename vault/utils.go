package vault

import "os"

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(filename string) bool {
	stat, _ := os.Stat(filename)
	if stat == nil {
		return false
	}
	return stat.IsDir()
}

func IsFile(filename string) bool {
	return !IsDir(filename)
}
