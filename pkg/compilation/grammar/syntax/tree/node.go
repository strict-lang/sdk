package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

// Node is implemented by every node of the tree.
type Node interface {
	Region() input.Region
	// Accept makes the visitor visit this node.
	Accept(visitor Visitor)
	// AcceptRecursive makes the visitor visit this node and its children.
	AcceptRecursive(visitor Visitor)
}

// Named is implemented by all nodes that have a name.
type Named interface {
	// Name returns the nodes name.
	Name() string
}
