package widgets

import (
	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-ascii"
)

type Password struct {
	Message   string
	Value     string
	Config    WidgetConfig
	Indicator string
}

func (c *Password) Run() {

	config := c.Config
	if config.Writer == nil {
		config = DefaultConfig
	}

	if c.Indicator == "" {
		c.Indicator = ascii.Bullet
	}

	writer := config.Writer

	cursor := ascii.Cursor{writer}

	write(writer, "%s ", config.MessageColor.Color(c.Message))

	//c.Theme.Printf("%s ", label)
	x := 0

	buffer := ""

	for {
		a, _, _ := ascii.GetChar()
		ascii.HandleSignals(a)
		if a == ascii.Backspace {
			if x == 0 {
				continue
			}

			write(writer, "\b \b")

			x--
			buffer = buffer[0:x]
			continue

		} else if a == ascii.Enter {
			c.Value = buffer
			break
		}

		buffer += string(a)

		write(writer, config.StdinColor.Color(c.Indicator))

		//write(writer, "len %d", len(c.Indicator))
		x++
	}

	cursor.Backward(x)
	str := ""
	for x > 0 {
		str += c.Indicator // "*"
		x--
	}
	write(writer, "%s\n", config.HighlightColor.Color(str))

}
