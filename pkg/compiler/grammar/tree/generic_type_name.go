package tree

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// TODO(merlinosayimwen) Change this to a slice to support
//  types like maps and tuples.

type GenericTypeName struct {
	Name    string
	Generic TypeName
	Region  input.Region
}

func (name *GenericTypeName) FullName() string {
	return fmt.Sprintf("%s<%s>", name.Name, name.Generic.FullName())
}

func (name *GenericTypeName) NonGenericName() string {
	return name.Name
}

func (name *GenericTypeName) Accept(visitor Visitor) {
	visitor.VisitGenericTypeName(name)
}

func (name *GenericTypeName) AcceptRecursive(visitor Visitor) {
	name.Accept(visitor)
	name.Generic.AcceptRecursive(visitor)
}

func (name *GenericTypeName) Locate() input.Region {
	return name.Region
}
