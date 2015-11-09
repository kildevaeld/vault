package ascii

import "io"

type Cursor struct {
	Writer io.Writer
}

func (c Cursor) Move(x, y int) Cursor {
	return c.writeString(CursorMove(x, y))

}

func (c Cursor) MoveToX(x int) Cursor {
	return c.writeString((CursorToX(x)))
}

func (c Cursor) Forward(x int) Cursor {
	return c.writeString(CursorForward(x))
}

func (c Cursor) Backward(x int) Cursor {
	return c.writeString(CursorBackward(x))
}

func (c Cursor) Up(y int) Cursor {
	return c.writeString(CursorUp(y))
}

func (c Cursor) Down(y int) Cursor {
	return c.writeString(CursorDown(y))
}

func (c Cursor) Hide() Cursor {
	return c.writeString(CursorHide)
}

func (c Cursor) Show() Cursor {
	return c.writeString(CursorShow)
}

func (c Cursor) EraseLines(count int) Cursor {
	return c.writeString(EraseLines(count))
}

func (c Cursor) writeString(str string) Cursor {
	c.Writer.Write([]byte(str))
	return c
}
