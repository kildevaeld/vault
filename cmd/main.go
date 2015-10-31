package main

import (
	"fmt"
	"os"

	"github.com/kildevaeld/vault"
)

var (
	//FILE = "worldcities.csv"
	FILE = "main.go"
)

func main() {

	stat, _ := os.Stat(FILE)

	size := stat.Size()
	pkgSize := int64(1024 * 32)

	fmt.Printf("File size %d - %d\n", size, size%pkgSize)
	//var bb []byte
	//var n [][]byte
	file, _ := os.Open(FILE)
	defer file.Close()
	dest, _ := os.Create("test")
	defer dest.Close()

	key := vault.Key([]byte("key"))
	err := vault.Encrypt(dest, file, stat.Size(), key)

	if err != nil {
		fmt.Printf("%v\n", err)
	}

	tStat, _ := os.Stat("test")

	testFile, _ := os.Open("test")
	defer testFile.Close()
	err = vault.Decrypt(os.Stdout, testFile, tStat.Size(), key)

	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
