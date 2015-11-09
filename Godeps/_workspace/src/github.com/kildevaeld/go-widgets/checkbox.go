package widgets

import (
	"strings"

	"github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-ascii"
)

type Checkbox struct {
	Message             string
	Choices             []string
	Value               []string
	SelectedIndicator   string
	UnselectedIndicator string
	Config              WidgetConfig
}

func (c *Checkbox) Run() {
	choices := c.Choices
	var results []string
	config := c.Config
	if config.Writer == nil {
		config = DefaultConfig
		c.Config = config
	}

	if c.SelectedIndicator == "" {
		c.SelectedIndicator = ascii.CircleCross
	}
	if c.UnselectedIndicator == "" {
		c.UnselectedIndicator = ascii.RadioOff
	}

	writer := config.Writer

	cursor := ascii.Cursor{writer}

	cursor.Hide()

	write(writer, "%s\n", config.MessageColor.Color(c.Message))

	for i, s := range choices {
		if i == len(choices)-1 {
			c.printLine(results, s, true)
		} else {
			c.printLine(results, s, false)
		}
		write(writer, "\n")
	}
	l := len(choices)

	cursor.Up(1)
	curPos := l - 1
	x := 0
	for {
		a, k, e := ascii.GetChar()
		if e != nil {
			return
		}

		ascii.HandleSignals(a)

		if k == ascii.UpKeyCode && curPos != 0 {
			cursor.Backward(x)
			x = c.printLine(results, choices[curPos], false)

			curPos = curPos - 1
			cursor.Up(1).Backward(x)

			x = c.printLine(results, choices[curPos], true)

		} else if k == ascii.DownKeyCode && curPos < l-1 {
			cursor.Backward(x)
			x = c.printLine(results, choices[curPos], false)

			curPos = curPos + 1
			cursor.Down(1).Backward(x)

			x = c.printLine(results, choices[curPos], true)

		} else if a == ascii.Enter {
			break
		} else if a == ascii.Space {
			cursor.Backward(x)

			if i := contains(results, choices[curPos]); i > -1 {
				results = append(results[:i], results[i+1:]...)
				x = c.printLine(results, choices[curPos], false)
			} else {
				results = append(results, choices[curPos])
				x = c.printLine(results, choices[curPos], true)
			}
		}
	}
	c.Value = results

	cursor.Down(l - curPos)

	for l > -1 {
		cursor.Up(1)
		write(writer, ascii.ClearLine)

		l = l - 1
	}
	vals := strings.Join(results, ", ")
	write(writer, "%s %s\n", config.MessageColor.Color(c.Message), config.HighlightColor.Color(vals))

	cursor.Show()
	return
}

func (c *Checkbox) printLine(results []string, s string, highlight bool) int {
	i := c.getIndicator(results, s)
	color := c.Config.StdinColor
	if highlight {
		color = c.Config.HighlightColor
	}

	return write(c.Config.Writer, color.Color("%s %s %s", ascii.EraseEndLine, i, s))
}
func (c *Checkbox) getIndicator(results []string, s string) string {
	i := c.UnselectedIndicator
	if contains(results, s) > -1 {
		i = c.SelectedIndicator
	}
	return i
}
