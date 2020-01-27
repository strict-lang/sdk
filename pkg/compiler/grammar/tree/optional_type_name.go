package tree

import "strict.dev/sdk/pkg/compiler/input"

type OptionalTypeName struct {
	Region input.Region
	TypeName TypeName
	Parent Node
}

func (name *OptionalTypeName) FullName() string {
	return name.TypeName.FullName() + "?"
}

func (name *OptionalTypeName) BaseName() string {
	return name.TypeName.BaseName()
}

func (name *OptionalTypeName) Matches(target Node) bool {
	if _, isWildcard := target.(*WildcardNode); isWildcard {
		return true
	}
	targetName, ok := target.(*OptionalTypeName)
	return ok && targetName.TypeName.Matches(name.TypeName)
}

func (name *OptionalTypeName) Locate() input.Region {
	return name.Region
}

func (name *OptionalTypeName) Accept(visitor Visitor) {
}

func (name *OptionalTypeName) AcceptRecursive(visitor Visitor) {
	name.TypeName.AcceptRecursive(visitor)
}

func (name *OptionalTypeName) SetEnclosingNode(node Node) {
	name.Parent = node
}

func (name *OptionalTypeName) EnclosingNode() (Node, bool) {
	return name.Parent, name.Parent != nil
}