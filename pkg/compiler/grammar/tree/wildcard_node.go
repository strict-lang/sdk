package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/scope"
)

// WildcardNode is a node that can be used as a leaf node for other nodes.
// It is only used for testing. Wildcards do not have children and are not
// created by the parser.
type WildcardNode struct {
	Statement
	Expression
	Region input.Region
	Parent Node
	scope  scope.Scope
}

func (node *WildcardNode) UpdateScope(target scope.Scope) {
	node.scope = target
}

func (node *WildcardNode) Scope() scope.Scope {
	return node.scope
}

func (node *WildcardNode) SetEnclosingNode(target Node) {
	node.Parent = target
}

func (node *WildcardNode) EnclosingNode() (Node, bool) {
	return node.Parent, node.Parent != nil
}

func (node *WildcardNode) Accept(visitor Visitor) {
	visitor.VisitWildcardNode(node)
}

func (node *WildcardNode) AcceptRecursive(visitor Visitor) {
	node.Accept(visitor)
}

func (node *WildcardNode) Locate() input.Region {
	return node.Region
}

func (node *WildcardNode) Matches(target Node) bool {
	return true
}
