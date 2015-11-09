package gocmd

import (
	"fmt"
	"strings"
	"reflect"
	"bufio"
	"os"
	"io"
)

type CmdServer interface {
	Register(interface{})
	Serve() int
	SetPrompt(prompt string)

	Exit(int)
}

type defaultCmdServer struct {
	prompt string
	cmdmap map[string]reflect.Value
	isQuit bool
	quitCode int
}

func NewServer(prompt string) CmdServer {
	return &defaultCmdServer{prompt, make(map[string]reflect.Value), false, 0}
}

func (s *defaultCmdServer) SetPrompt(prompt string) {
	s.prompt = prompt
}

func (s *defaultCmdServer) Register(service interface{}) {
	t := reflect.TypeOf(service)
	v := reflect.ValueOf(service)
	for i, n := 0, t.NumMethod(); i<n; i++ {
		m := t.Method(i)

		if m.PkgPath != "" {
			// Method must be exported.
			continue
		}

		if m.Type.NumOut() != 0 {
			continue
		}

		if m.Type.NumIn() != 3 {
			continue
		}

		argsType := m.Type.In(2)
		if argsType.Kind() != reflect.Slice {
			continue
		}
		if argsType.Elem().Kind() != reflect.String {
			continue
		}
		s.cmdmap[strings.ToLower(m.Name)] = v.Method(i)
	}
	return
}

func (s *defaultCmdServer) Serve() int {
	reader := bufio.NewReader(os.Stdin)
	for !s.isQuit {
		fmt.Printf("%s", s.prompt)
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			s.quitCode = -1
			break
		}
		input = strings.TrimRight(input, "\n \t")
		args := strings.Split(input, " ")
		if method, ok := s.cmdmap[args[0]]; ok {
			var cmd CmdServer = s
			callArgs := []reflect.Value{reflect.ValueOf(cmd), reflect.ValueOf(args[1:])}
			method.Call(callArgs)
		} else {
			if len(args[0]) > 0 {
				fmt.Printf("Bad command: %s\n", args[0])
			}
		}
	}
	return s.quitCode
}

func (s *defaultCmdServer) Exit(code int) {
	s.isQuit = true
	s.quitCode = code
}
