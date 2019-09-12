package ast

import (
	"gitlab.com/strict-lang/sdk/compilation/token"
)

type MethodCall struct {
	// Method is the called method. It can be any kind of expression
	// with the value of a method. Common nodes are identifiers and
	// field selectors.
	Method Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments    []Node
	NodePosition Position
}

func (call *MethodCall) Accept(visitor *Visitor) {
	visitor.VisitMethodCall(call)
}

func (call *MethodCall) AcceptRecursive(visitor *Visitor) {
	visitor.VisitMethodCall(call)
	call.Method.AcceptRecursive(visitor)
	for _, argument := range call.Arguments {
		argument.AcceptRecursive(visitor)
	}
}

func (call *MethodCall) Position() Position {
	return call.Position()
}

type Identifier struct {
	Value        string
	NodePosition Position
}

func (identifier *Identifier) Accept(visitor *Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) AcceptRecursive(visitor *Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) Position() Position {
	return identifier.Position()
}

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	Operator     token.Operator
	Operand      Node
	NodePosition Position
}

func (unary *UnaryExpression) Accept(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
}

func (unary *UnaryExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
	unary.Operand.AcceptRecursive(visitor)
}

func (unary *UnaryExpression) Position() Position {
	return unary.Position()
}

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Node
	RightOperand Node
	Operator     token.Operator
	NodePosition Position
}

func (binary *BinaryExpression) Accept(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
}

func (binary *BinaryExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
	binary.LeftOperand.AcceptRecursive(visitor)
	binary.RightOperand.AcceptRecursive(visitor)
}

func (binary *BinaryExpression) Position() Position {
	return binary.NodePosition
}

type SelectExpression struct {
	Target       Node
	Selection    Node
	NodePosition Position
}

func (expression *SelectExpression) Accept(visitor *Visitor) {
	visitor.VisitSelectorExpression(expression)
}

func (expression *SelectExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitSelectorExpression(expression)
	expression.Target.AcceptRecursive(visitor)
	expression.Selection.AcceptRecursive(visitor)
}

func (expression *SelectExpression) Position() Position {
	return expression.NodePosition
}

type CreateExpression struct {
	NodePosition Position
	Constructor  *MethodCall
	Type         TypeName
}

func (create *CreateExpression) Accept(visitor *Visitor) {
	visitor.VisitCreateExpression(create)
}

func (create *CreateExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitCreateExpression(create)
	create.Constructor.AcceptRecursive(visitor)
}

func (create *CreateExpression) Position() Position {
	return create.NodePosition
}

type ListSelectExpression struct {
	Index Node
	Target Node
	NodePosition Position
}

func (expression *ListSelectExpression) Accept(visitor *Visitor) {
	visitor.VisitListSelectExpression(expression)
}

func (expression *ListSelectExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitListSelectExpression(expression)
	expression.Index.AcceptRecursive(visitor)
	expression.Target.AcceptRecursive(visitor)
}

func (expression *ListSelectExpression) Position() Position {
	return expression.NodePosition
}