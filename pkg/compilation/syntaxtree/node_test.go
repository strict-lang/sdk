package syntaxtree

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type TestPosition struct {
	BeginOffset source2.Offset
	EndOffset   source2.Offset
}

func (position *TestPosition) Begin() source2.Offset {
	return position.BeginOffset
}

func (position *TestPosition) End() source2.Offset {
	return position.EndOffset
}
