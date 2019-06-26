package ast

import "fmt"

type Statement interface {
	Node
	statement()
}

type ExpressionStatement struct {
	expression Expression
}

func (expression ExpressionStatement) Expression() Expression {
	return expression.expression
}

func (expression ExpressionStatement) statement() { }

func (expression ExpressionStatement) Accept(visitor Visitor) {
	visitor(expression.expression)
}

type MethodCall struct {
	Typed
	Method *Member
	Arguments []Expression
}

func (call MethodCall) Type() *Type {
	return call.Method.ValueType
}

func (call MethodCall) statement() { }

func (call MethodCall) String() string {
	if call.Method == nil {
		return fmt.Sprintf("AnynomousMethodCall(%s)", call.Arguments)
	}
	return fmt.Sprintf("{%s %s(%s)}",
			call.Type(), call.Method.Name, call.Arguments)
}

type BlockStatement struct {

}

type ConditionalStatement struct {
	Condition Expression
	Body      *BlockStatement
	Else      Statement
	position Position
}

func (conditional ConditionalStatement) Position() Position {
	return conditional.position
}

type ForLoopStatement struct {
	Initialization Expression
	Termination Expression
	Increment Expression
	Block BlockStatement
	position Position
}

func (forLoop ForLoopStatement) Position() Position {
	return forLoop.position
}

type PreIncrementStatement struct {

}

type PostIncrementStatement struct {}