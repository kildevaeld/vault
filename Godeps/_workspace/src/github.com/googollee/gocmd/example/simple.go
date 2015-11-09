package main

import (
	"fmt"
	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/googollee/gocmd"
	"strings"
)

type TestServer struct {
}

func (s *TestServer) Ls(server gocmd.CmdServer, args []string) {
	fmt.Println("ls")
	for a := range args {
		fmt.Println(a)
	}
	fmt.Printf("Total %d.\n", len(args))
}

func (s *TestServer) Echo(server gocmd.CmdServer, args []string) {
	fmt.Println(strings.Join(args, " "))
}

func (s *TestServer) Quit(server gocmd.CmdServer, args []string) {
	server.Exit(0)
}

func main() {
	s := gocmd.NewServer(">")
	s.Register(new(TestServer))
	s.Serve()
}
