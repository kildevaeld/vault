package widgets

import "github.com/kildevaeld/vault/Godeps/_workspace/src/github.com/kildevaeld/go-ascii"

type Paginate int

const (
	PaginateBack Paginate = 1 << iota
	PaginateNext
	PaginateDone
)

type PaginatedList struct {
	Message   string
	Paginate  func(page int) []string
	Config    WidgetConfig
	Indicator string
	cursor    ascii.Cursor
	Value     string
}

func (p *PaginatedList) printList(choices []string, result *string) Paginate {

	for i, choice := range choices {
		if i == len(choices)-1 {
			p.highlightLine(choice)
		} else {
			p.printLine(choice)
		}

	}

	x := p.writeString("p(rev), n(ext) or d(one)\r")
	p.cursor.Up(1)
	x = 0
	l := len(choices)
	y := l - 1
	height := len(choices)
	for {
		c, j, _ := ascii.GetChar()
		ascii.HandleSignals(c)
		if c == 'n' {
			p.cursor.Down(l - y).EraseLines(height)
			return PaginateNext
		} else if c == 'p' {
			p.cursor.Down(l - y).EraseLines(height)
			return PaginateBack
		} else if c == 'd' {
			p.cursor.Down(l - y).EraseLines(height)
			return PaginateDone
		} else if j == ascii.UpKeyCode && y != 0 {
			p.cursor.Backward(x)
			x = p.printLine(choices[y])

			y -= 1
			p.cursor.Up(2).Backward(x)
			x = p.highlightLine(choices[y])
			p.cursor.Up(1)
		} else if j == ascii.DownKeyCode && y < l-1 {

			p.cursor.Backward(x)
			x = p.printLine(choices[y])

			y += 1
			p.cursor.Backward(x)
			x = p.highlightLine(choices[y])
			p.cursor.Up(1)
		} else if c == ascii.Enter {
			p.cursor.Down(l - y).EraseLines(height)
			*result = choices[y]
			break
		}
	}

	return PaginateDone

}

func (p *PaginatedList) highlightLine(str string) int {
	return p.writeString(p.Config.HighlightColor.Color(" %s %s\n", p.Indicator, str))
}

func (p *PaginatedList) printLine(str string) int {
	return p.writeString(p.Config.StdinColor.Color("   %s\n", str))

}

func (p *PaginatedList) writeString(str string, args ...interface{}) int {
	return write(p.Config.Writer, str, args...)
}

func (p *PaginatedList) Run() {
	config := p.Config
	if config.Writer == nil {
		config = DefaultConfig
		p.Config = config
	}

	if p.Indicator == "" {
		p.Indicator = ascii.Pointer
	}

	p.cursor = ascii.Cursor{p.Config.Writer}

	page := 1
	var result string
	p.cursor.Hide()
	p.writeString("%s\n", p.Message)

	var choices []string
	for {

		if page >= 1 {
			choices = p.Paginate(page)
		} else {
			page = 1
		}

		if choices == nil {
			if page > 1 {
				page = 1
				continue
			}
			break
		}

		action := p.printList(choices, &result)
		p.cursor.Up(1)
		if action == PaginateNext {
			page = page + 1
		} else if action == PaginateBack {
			page = page - 1
		} else {
			break
		}
		//ui.Printf("%s", ascii.EraseLines(12))
		//p.cursor.Backward(25)
	}

	p.cursor.Up(1)
	p.writeString("%s %s\n", p.Message, p.Config.HighlightColor.Color(result))
	p.Value = result
	p.cursor.Show()
}

func NewPaginatedList(msg string, fn func(page int) []string) *PaginatedList {
	return &PaginatedList{
		Message:  msg,
		Paginate: fn,
		Config:   DefaultConfig,
	}
}
