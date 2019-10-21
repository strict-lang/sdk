package tree

type TypeName interface {
	Node
	FullName() string
	NonGenericName() string
}
