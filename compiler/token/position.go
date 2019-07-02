package token

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/source"
)

type Position struct {
	Begin source.Offset
	End   source.Offset
}

func (position Position) String() string {
	return fmt.Sprintf("Position{%d..%d}", position.Begin, position.End)
}
func (position Position) Contains(offset source.Offset) bool {
	return position.Begin <= offset && offset <= position.End
}

func (position Position) ContainsPosition(target Position) bool {
	return position.Contains(target.Begin) && position.Contains(target.End)
}
