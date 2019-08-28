package ast

import (
	"gitlab.com/strict-lang/sdk/compiler/token"
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

func (call *MethodCall) AcceptAll(visitor *Visitor) {
	visitor.VisitMethodCall(call)
	call.Method.AcceptAll(visitor)
	for _, argument := range call.Arguments {
		argument.AcceptAll(visitor)
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

func (identifier *Identifier) AcceptAll(visitor *Visitor) {
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

func (unary *UnaryExpression) AcceptAll(visitor *Visitor) {
	visitor.VisitUnaryExpression(unary)
	unary.Operand.AcceptAll(visitor)
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

func (binary *BinaryExpression) AcceptAll(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
	binary.LeftOperand.AcceptAll(visitor)
	binary.RightOperand.AcceptAll(visitor)
}

func (binary *BinaryExpression) Position() Position {
	return binary.NodePosition
}

type SelectorExpression struct {
	Target       Node
	Selection    Node
	NodePosition Position
}

func (selector *SelectorExpression) Accept(visitor *Visitor) {
	visitor.VisitSelectorExpression(selector)
}

func (selector *SelectorExpression) AcceptAll(visitor *Visitor) {
	visitor.VisitSelectorExpression(selector)
	selector.Target.AcceptAll(visitor)
	selector.Selection.AcceptAll(visitor)
}

func (selector *SelectorExpression) Position() Position {
	return selector.NodePosition
}

type CreateExpression struct {
	NodePosition Position
	Constructor  *MethodCall
}

func (create *CreateExpression) Accept(visitor *Visitor) {
	visitor.VisitCreateExpression(create)
}

func (create *CreateExpression) AcceptAll(visitor *Visitor) {
	visitor.VisitCreateExpression(create)
	create.Constructor.AcceptAll(visitor)
}

func (create *CreateExpression) Position() Position {
	return create.NodePosition
}
