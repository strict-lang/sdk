package tree

type Statement interface {
	Node
	IsModifyingControlFlow() bool
}

