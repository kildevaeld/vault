package ascii

import "fmt"

type Style struct {
	open  int
	close int
}

func (s Style) Open() string {
	return fmt.Sprintf("%s%dm", ESC, s.open)
}

func (s Style) Close() string {
	return fmt.Sprintf("%s%dm", ESC, s.close)
}

func (s Style) Color(msg string, args ...interface{}) string {
	return fmt.Sprintf("%s%s%s", s.Open(), fmt.Sprintf(msg, args...), s.Close())
}

var (
	// Modifier
	Reset        = Style{0, 0}
	Bold         = Style{1, 22}
	Dim          = Style{2, 22}
	Italic       = Style{3, 23}
	Underline    = Style{4, 24}
	Inverse      = Style{7, 27}
	Hidden       = Style{8, 28}
	Strikethough = Style{9, 29}
	// Colors
	Black   = Style{30, 39}
	Red     = Style{31, 39}
	Green   = Style{32, 39}
	Yellow  = Style{33, 39}
	Blue    = Style{34, 39}
	Megenta = Style{35, 39}
	Cyan    = Style{36, 39}
	White   = Style{37, 39}
	Gray    = Style{90, 39}
	// Background Colors
	BgBlack   = Style{40, 49}
	BgRed     = Style{41, 49}
	BgGreen   = Style{42, 49}
	BgYellow  = Style{43, 49}
	BgBlue    = Style{44, 49}
	BgMegenta = Style{45, 49}
	BgCyan    = Style{46, 49}
	BgWhite   = Style{47, 49}
)
