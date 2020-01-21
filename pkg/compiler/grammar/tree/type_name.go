package tree

type TypeName interface {
	Node
	// FullName returns the full type name, including generics and punctuation.
	FullName() string
	// BaseName returns base type names, without punctuation and generics.
	// Example: Number[] -> Number, MutableList<String> -> MutableList.
	BaseName() string
}
