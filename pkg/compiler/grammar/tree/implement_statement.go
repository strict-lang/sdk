package tree

import "strict.dev/sdk/pkg/compiler/input"

type ImplementStatement struct {
	Parent Node
	Region input.Region
	Trait  TypeName
}

func (statement *ImplementStatement) Accept(visitor Visitor) {
}

func (statement *ImplementStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	statement.Trait.AcceptRecursive(visitor)
}

func (statement *ImplementStatement) Locate() input.Region {
	return statement.Region
}

func (statement *ImplementStatement) Matches(target Node) bool {
	if _, isWildcard := target.(*WildcardNode); isWildcard {
		return true
	}
	targetStatement, ok := target.(*ImplementStatement)
	return ok && statement.Trait.Matches(targetStatement.Trait)
}

func (statement *ImplementStatement) SetEnclosingNode(node Node) {
	statement.Parent = node
}

func (statement *ImplementStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}
