package syntaxtree

import (
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

type CallExpression struct {
	// Method is the called method. It can be any kind of expression
	// with the value of a method. Common nodes are identifiers and
	// field selectors.
	Method Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments    []*CallArgument
	NodePosition Position
}

func (call *CallExpression) Accept(visitor *Visitor) {
	visitor.VisitCallExpression(call)
}

func (call *CallExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitCallExpression(call)
	AcceptRecursive(visitor)
	for _, argument := range call.Arguments {
		argument.AcceptRecursive(visitor)
	}
}

func (call *CallExpression) Position() Position {
	return call.Position()
}

type CallArgument struct {
	Label        string
	Value        Node
	NodePosition Position
}

func (argument *CallArgument) IsLabeled() bool {
	return argument.Label != ""
}

func (argument *CallArgument) Accept(visitor *Visitor) {
	visitor.VisitCallArgument(argument)
}

func (argument *CallArgument) AcceptRecursive(visitor *Visitor) {
	visitor.VisitCallArgument(argument)
	AcceptRecursive(visitor)
}

func (argument *CallArgument) Position() Position {
	return argument.NodePosition
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
	Operator     token2.Operator
	Operand      Node
	NodePosition Position
}

func (unary *UnaryExpression) Accept(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
}

func (unary *UnaryExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
	AcceptRecursive(visitor)
}

func (unary *UnaryExpression) Position() Position {
	return unary.Position()
}

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Node
	RightOperand Node
	Operator     token2.Operator
	NodePosition Position
}

func (binary *BinaryExpression) Accept(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
}

func (binary *BinaryExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
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
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (expression *SelectExpression) Position() Position {
	return expression.NodePosition
}

type CreateExpression struct {
	NodePosition Position
	Constructor  *CallExpression
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
	Index        Node
	Target       Node
	NodePosition Position
}

func (expression *ListSelectExpression) Accept(visitor *Visitor) {
	visitor.VisitListSelectExpression(expression)
}

func (expression *ListSelectExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitListSelectExpression(expression)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (expression *ListSelectExpression) Position() Position {
	return expression.NodePosition
}
