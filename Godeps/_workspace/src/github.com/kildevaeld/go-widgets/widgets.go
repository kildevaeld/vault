package widgets

import (
	"fmt"
	"io"
	"os"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-ascii"
)

func contains(haystack []string, needle string) int {
	for i, n := range haystack {
		if needle == n {
			return i
		}
	}
	return -1
}

type WidgetConfig struct {
	MessageColor   ascii.Style
	HighlightColor ascii.Style
	StdinColor     ascii.Style
	Writer         io.Writer
}

var DefaultConfig WidgetConfig = WidgetConfig{
	MessageColor:   ascii.Dim,
	HighlightColor: ascii.Cyan,
	StdinColor:     ascii.Reset,
	Writer:         os.Stdout,
}

type Widget interface {
	Run()
}

func write(w io.Writer, msg string, args ...interface{}) int {
	str := fmt.Sprintf(msg, args...)
	w.Write([]byte(str))
	return len(str)
}
