package syntaxtree

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type TestPosition struct {
	BeginOffset source.Offset
	EndOffset   source.Offset
}

func (position *TestPosition) Begin() source.Offset {
	return position.BeginOffset
}

func (position *TestPosition) End() source.Offset {
	return position.EndOffset
}
