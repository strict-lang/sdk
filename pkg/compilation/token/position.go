package token

import (
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type Position struct {
	BeginOffset source.Offset
	EndOffset   source.Offset
}

func (position Position) Begin() source.Offset {
	return position.BeginOffset
}

func (position Position) End() source.Offset {
	return position.EndOffset
}

func (position Position) String() string {
	return fmt.Sprintf("Position{%d..%d}", position.BeginOffset, position.EndOffset)
}
func (position Position) Contains(offset source.Offset) bool {
	return position.BeginOffset <= offset && offset <= position.EndOffset
}

func (position Position) ContainsPosition(target Position) bool {
	return position.Contains(target.BeginOffset) && position.Contains(target.EndOffset)
}
