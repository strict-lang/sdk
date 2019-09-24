package syntaxtree

var (
	_ Node = &Identifier{}
	_ Node = &SelectExpression{}
	_ Node = &BinaryExpression{}
	_ Node = &UnaryExpression{}
	_ Node = &CallExpression{}
	_ Node = &CreateExpression{}
)
