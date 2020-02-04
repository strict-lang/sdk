package tree

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/input"
)

type ListTypeName struct {
	Element       TypeName
	Region        input.Region
	Parent        Node
	typeReference *TypeReference
}

func (name *ListTypeName) TypeReference() *TypeReference {
	return name.typeReference
}

func (name *ListTypeName) SetEnclosingNode(target Node) {
	name.Parent = target
}

func (name *ListTypeName) EnclosingNode() (Node, bool) {
	return name.Parent, name.Parent != nil
}

func (name *ListTypeName) FullName() string {
	return fmt.Sprintf("%s[]", name.Element.FullName())
}

func (name *ListTypeName) BaseName() string {
	return name.Element.BaseName()
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

func (name *ListTypeName) Matches(node Node) bool {
	if target, ok := node.(*ListTypeName); ok {
		return name.Element.Matches(target.Element)
	}
	return false
}
