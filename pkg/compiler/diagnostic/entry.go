package diagnostic

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
)

type Entry struct {
	Kind     *Kind
	Stage    *Stage
	Source   string
	Message  string
	UnitName string
	Position Position
	Error    *RichError
}

type Position struct {
	Begin input.Position
	End   input.Position
}

func (position Position) isSpanningMultipleLines() bool {
	return position.Begin.Line.Index != position.End.Line.Index
}

func (position Position) endColumn() input.Offset {
	if position.isSpanningMultipleLines() {
		return position.Begin.Line.Length
	}
	return position.End.Column
}
