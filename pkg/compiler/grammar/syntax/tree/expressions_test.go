package tree

var (
	_ Node = &Identifier{}
	_ Node = &FieldSelectExpression{}
	_ Node = &BinaryExpression{}
	_ Node = &UnaryExpression{}
	_ Node = &CallExpression{}
	_ Node = &CreateExpression{}
)
