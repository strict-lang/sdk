package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

// WildcardNode is a node that can be used as a leaf node for other nodes.
// It is only used for testing. Wildcards do not have children and are not
// created by the parser.
type WildcardNode struct {
	Statement
	Expression
	Region input.Region
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
