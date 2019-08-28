package ast

import "gitlab.com/strict-lang/sdk/compilation/source"

type TestPosition struct {
	BeginOffset source.Offset
	EndOffset source.Offset
}

var ZeroPosition Position = &TestPosition{
	BeginOffset: 0,
	EndOffset: 0,
}

func (position *TestPosition) Begin() source.Offset {
	return position.BeginOffset
}

func (position *TestPosition) End() source.Offset {
	return position.EndOffset
}