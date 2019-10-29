package tree

type TypeDescriptor string

type resolvedType struct {
	descriptor TypeDescriptor
	isResolved bool
}

func (resolvedType *resolvedType) setDescriptor(descriptor TypeDescriptor) {
	resolvedType.descriptor = descriptor
	resolvedType.isResolved = true
}

func (resolvedType *resolvedType) getDescriptor() (TypeDescriptor, bool) {
	return resolvedType.descriptor, resolvedType.isResolved
}
