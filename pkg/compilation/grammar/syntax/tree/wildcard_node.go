package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

type WildcardNode struct {
	NodeRegion input.Region
}

func (node* WildcardNode) Accept(visitor Visitor) {
	visitor.VisitWildcardNode(node)
}

func (node* WildcardNode) AcceptRecursive(visitor Visitor) {
	node.Accept(visitor)
}

func (node *WildcardNode) Region() input.Region {
	return node.NodeRegion
}
