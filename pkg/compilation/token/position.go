package token

import (
	"fmt"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type Position struct {
	BeginOffset source2.Offset
	EndOffset   source2.Offset
}

func (position Position) Begin() source2.Offset {
	return position.BeginOffset
}

func (position Position) End() source2.Offset {
	return position.EndOffset
}

func (position Position) String() string {
	return fmt.Sprintf("Position{%d..%d}", position.BeginOffset, position.EndOffset)
}
func (position Position) Contains(offset source2.Offset) bool {
	return position.BeginOffset <= offset && offset <= position.EndOffset
}

func (position Position) ContainsPosition(target Position) bool {
	return position.Contains(target.BeginOffset) && position.Contains(target.EndOffset)
}
