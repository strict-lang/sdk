package syntaxtree

// Node is implemented by every node of the syntaxtree.
type Node interface {
	Positioned

	// Accept makes the visitor visit this node.
	Accept(visitor *Visitor)
	// AcceptRecursive makes the visitor visit this node and its children.
	AcceptRecursive(visitor *Visitor)
}

// Named is implemented by all nodes that have a name.
type Named interface {
	// Name returns the nodes name.
	Name() string
}
