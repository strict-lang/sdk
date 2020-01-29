package tree

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/input"
)

// TODO(merlinosayimwen) Change this to a slice to support
//  types like maps and tuples.

type GenericTypeName struct {
	Name    string
	Generic TypeName
	Region  input.Region
	Parent Node
	typeReference *TypeReference
}

func (name *GenericTypeName) TypeReference() *TypeReference {
	return name.typeReference
}

func (name *GenericTypeName) SetEnclosingNode(target Node) {
  name.Parent = target
}

func (name *GenericTypeName) EnclosingNode() (Node, bool) {
  return name.Parent, name.Parent != nil
}

func (name *GenericTypeName) FullName() string {
	return fmt.Sprintf("%s<%s>", name.Name, name.Generic.FullName())
}

func (name *GenericTypeName) BaseName() string {
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

func (name *GenericTypeName) Matches(node Node) bool {
	if target, ok := node.(*GenericTypeName); ok {
		return name.Name == target.Name &&
			name.Generic.Matches(target.Generic)
	}
	return false
}
