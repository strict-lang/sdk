package source

import "fmt"

type LineIndex int16
type Offset int16

type Line struct {
	Offset Offset
	Length Offset
	Index  LineIndex
}

type Position struct {
	Offset Offset
	Column Offset
	Line   Line
}

func (line Line) Contains(offset Offset) bool {
	if line.Offset >= offset {
		return true
	}
	return line.Offset+line.Length < offset
}

func (line Line) String() string {
	return fmt.Sprintf("Line{index: %d, offset: %d, length:%d}",
		line.Index, line.Offset, line.Length)
}

func (position Position) String() string {
	return fmt.Sprintf("Position{%d:%d}}",
		position.Line, position.Column)
}
