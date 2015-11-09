package main

import (
	"fmt"

	ascii "github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-ascii"
)

func main() {

	s := ascii.CursorShow

	fmt.Printf("%sLest%s\n", ascii.Red.Open(), ascii.Red.Close())

	fmt.Printf(s)

	fmt.Scanf("%s")

}
