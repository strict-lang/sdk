package tree

type Expression interface {
	Node
	// ResolveType(descriptor TypeDescriptor)
	// GetResolvedType() (TypeDescriptor, bool)
}

type StoredExpression interface {
	Expression
}
