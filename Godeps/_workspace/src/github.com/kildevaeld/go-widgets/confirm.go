package widgets

import (
	"time"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-ascii"
)

type Confirm struct {
	Message string
	Value   bool
	Config  WidgetConfig
}

func (c *Confirm) Run() {

	config := c.Config
	if config.Writer == nil {
		config = DefaultConfig
	}
	writer := config.Writer

	write(writer, config.MessageColor.Color("%s [yn] ", c.Message))

	a, _, _ := ascii.GetChar()

	ascii.HandleSignals(a)

	ans := string(a)
	if ans == "y" || ans == "Y" {
		c.Value = true
		ans = "yes"
	} else if ans == "n" || ans == "N" {
		c.Value = false
		ans = "no"
	} else {
		write(writer, "%s%s ", ascii.ClearLine, c.Message)
		write(writer, config.HighlightColor.Color("please enter %s(es) or %s(o)", ascii.Bold.Color("y"), ascii.Bold.Color("n")))

		time.Sleep(1 * time.Second)
		write(writer, ascii.ClearLine)
		c.Run()
		return
	}

	write(writer, "%s\n", config.HighlightColor.Color(ans))

}
