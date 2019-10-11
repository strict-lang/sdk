package tree

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type ListTypeName struct {
	Element TypeName
	Region input.Region
}

func (name *ListTypeName) FullName() string {
	return fmt.Sprintf("%s[]", name.Element.FullName())
}

func (name *ListTypeName) NonGenericName() string {
	return name.Element.NonGenericName()
}

func (name *ListTypeName) Accept(visitor Visitor) {
	visitor.VisitListTypeName(name)
}

func (name *ListTypeName) AcceptRecursive(visitor Visitor) {
	name.Accept(visitor)
	name.Element.AcceptRecursive(visitor)
}

func (name *ListTypeName) Locate() input.Region {
	return name.Region
}
