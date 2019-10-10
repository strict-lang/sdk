package tree

var (
	_ Node = &StringLiteral{}
	_ Node = &NumberLiteral{}
)
