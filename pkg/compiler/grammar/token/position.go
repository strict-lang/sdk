package token

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/input"
)

type Position struct {
	BeginOffset input.Offset
	EndOffset   input.Offset
}

func (position Position) Begin() input.Offset {
	return position.BeginOffset
}

func (position Position) End() input.Offset {
	return position.EndOffset
}

func (position Position) String() string {
	return fmt.Sprintf("Area{%d..%d}", position.BeginOffset, position.EndOffset)
}
func (position Position) Contains(offset input.Offset) bool {
	return position.BeginOffset <= offset && offset <= position.EndOffset
}

func (position Position) ContainsPosition(target Position) bool {
	return position.Contains(target.BeginOffset) && position.Contains(target.EndOffset)
}
