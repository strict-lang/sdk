package tree

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/input"
)

type ListTypeName struct {
	Element TypeName
	Region  input.Region
}

func (name *ListTypeName) FullName() string {
	return fmt.Sprintf("%s[]", FullName())
}

func (name *ListTypeName) NonGenericName() string {
	return NonGenericName()
}

func (name *ListTypeName) Accept(visitor Visitor) {
	VisitListTypeName(name)
}

func (name *ListTypeName) AcceptRecursive(visitor Visitor) {
	name.Accept(visitor)
	AcceptRecursive(visitor)
}

func (name *ListTypeName) Mangle() string {
	return "ARRAY_" + Mangle()
}

func (name *ListTypeName) Locate() input.Region {
	return name.Region
}

func (name *ListTypeName) Matches(node Node) bool {
	if target, ok := node.(*ListTypeName); ok {
		return Matches(target.Element)
	}
	return false
}
