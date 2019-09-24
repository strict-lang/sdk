package syntaxtree

var (
	_ Node = &StringLiteral{}
	_ Node = &NumberLiteral{}
)
