package source

import "fmt"

type LineIndex int16
type Offset int16

type Position struct {
	Offset     Offset
	LineColumn Offset
	LineIndex  LineIndex
}

func (position Position) String() string {
	return fmt.Sprintf("Position { Line: %d, Column: %d }",
		position.LineIndex, position.LineColumn)
}
