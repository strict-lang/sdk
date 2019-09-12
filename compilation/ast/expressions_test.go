package ast

var (
	_ Node = &Identifier{}
	_ Node = &SelectExpression{}
	_ Node = &BinaryExpression{}
	_ Node = &UnaryExpression{}
	_ Node = &MethodCall{}
	_ Node = &CreateExpression{}
)
