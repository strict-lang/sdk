package tree

type Expression interface {
	Node
	// SetResolvedType(descriptor TypeDescriptor)
	// ResolvedType() (TypeDescriptor, bool)
}

type StoredExpression interface {
	Expression
}
