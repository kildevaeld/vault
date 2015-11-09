Gocmd: simple cmdline enviornment package
=========================================

Usage
-----

Command function must have signature like this:

        func (...) Command(server gocmd.CmdServer, args []string)

Examples can be found under example directory. Simple usage like below:

        package main

        import (
        	"github.com/googollee/gocmd"
        	"strings"
        	"fmt"
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
