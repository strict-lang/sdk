package tree

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/input"
)

type TestPosition struct {
	BeginOffset input.Offset
	EndOffset   input.Offset
}

func (position *TestPosition) Begin() input.Offset {
	return position.BeginOffset
}

func (position *TestPosition) End() input.Offset {
	return position.EndOffset
}
