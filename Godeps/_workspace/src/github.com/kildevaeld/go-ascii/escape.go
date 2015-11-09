package ascii

import (
	"fmt"
	"math"
)

const ESC = "\u001b["

// Cursor
const CursorLeft = ESC + "1000D"
const CursorSavePosition = ESC + "s"
const CursorRestorePosition = ESC + "u"
const CursorGetPosition = ESC + "6n"
const CursorNextLine = ESC + "E"
const CursorPrevLine = ESC + "F"
const CursorHide = ESC + "?25l"
const CursorShow = ESC + "?25h"

const EraseEndLine = ESC + "K"
const EraseStartLine = ESC + "1K"
const EraseLine = ESC + "2K"
const EraseDown = ESC + "J"
const EraseUp = ESC + "1J"
const EraseScreen = ESC + "2J"
const ScrollUp = ESC + "S"
const ScrollDown = ESC + "T"

const Beep = "\u0007"

func CursorTo(x, y int) string {
	return fmt.Sprintf("%s%d;%dH", ESC, y+1, x+1)
}

func CursorToX(x int) string {
	return fmt.Sprintf("%s%sG", ESC, (x + 1))
}

func CursorMove(x, y int) string {
	str := ""

	if x < 0 {
		str = fmt.Sprintf("%s%dD", ESC, math.Abs(float64(x)))
	} else {
		str = fmt.Sprintf("%s%dC", ESC, x)
	}

	if y < 0 {
		str = fmt.Sprintf("%s%s%dA", str, ESC, math.Abs(float64(y)))
	} else {
		str = fmt.Sprintf("%s%s%dB", str, ESC, y)
	}
	return str
}

func CursorUp(count int) string {
	return fmt.Sprintf("%s%dA", ESC, count)
}

func CursorDown(count int) string {
	return fmt.Sprintf("%s%dB", ESC, count)
}

func CursorForward(count int) string {
	return fmt.Sprintf("%s%dC", ESC, count)
}

func CursorBackward(count int) string {
	return fmt.Sprintf("%s%dD", ESC, count)
}

func EraseLines(count int) string {
	clear := ""

	for i := 0; i < count; i++ {
		postfix := ""
		if i < (count - 1) {
			postfix = CursorUp(1)
		}
		clear += CursorLeft + EraseEndLine + postfix
	}
	return clear

}
