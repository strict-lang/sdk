package ast

var (
	_ Node = &ExpressionStatement{}
	_ Node = &ForEachLoopStatement{}
	_ Node = &RangedLoopStatement{}
	_ Node = &ConditionalStatement{}
	_ Node = &IncrementStatement{}
	_ Node = &DecrementStatement{}
	_ Node = &InvalidStatement{}
	_ Node = &EmptyStatement{}
	_ Node = &TestStatement{}
	_ Node = &AssertStatement{}
	_ Node = &YieldStatement{}
	_ Node = &ReturnStatement{}
)
