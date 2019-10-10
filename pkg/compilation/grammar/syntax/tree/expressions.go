package tree

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
)

type Expression interface {
	Node
}

type StoredExpression interface {
	Expression
}

type CallExpression struct {
	// Target is the procedure that is called.
	Target Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments    []*CallArgument
	NodePosition InputRegion
}

func (call *CallExpression) Accept(visitor *Visitor) {
	visitor.VisitCallExpression(call)
}

func (call *CallExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitCallExpression(call)
	call.Target.AcceptRecursive(visitor)
	for _, argument := range call.Arguments {
		argument.AcceptRecursive(visitor)
	}
}

func (call *CallExpression) Area() InputRegion {
	return call.Area()
}

type CallArgument struct {
	Label        string
	Value        Node
	NodePosition InputRegion
}

func (argument *CallArgument) IsLabeled() bool {
	return argument.Label != ""
}

func (argument *CallArgument) Accept(visitor Visitor) {
	visitor.VisitCallArgument(argument)
}

func (argument *CallArgument) AcceptRecursive(visitor Visitor) {
	visitor.VisitCallArgument(argument)
	argument.Value.AcceptRecursive(visitor)
}

func (argument *CallArgument) Area()  input.Region{
	return argument.NodePosition
}

type Identifier struct {
	Value        string
	Region input.Region
}

func (identifier *Identifier) Accept(visitor Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) AcceptRecursive(visitor Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) Area() InputRegion {
	return identifier.Area()
}

// UnaryExpression is an operation on a single operand.
type UnaryExpression struct {
	Operator     token.Operator
	Operand      Expression
	NodeRegion   input.Region
}

func (unary *UnaryExpression) Accept(visitor Visitor) {
	visitor.VisitUnaryExpression(unary)
}

func (unary *UnaryExpression) AcceptRecursive(visitor Visitor) {
	unary.Accept(visitor)
	unary.Operand.AcceptRecursive(visitor)
}

func (unary *UnaryExpression) Locate() input.Region {
	return unary.NodeRegion
}

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Node
	RightOperand Node
	Operator     token.Operator
	NodePosition InputRegion
}

func (binary *BinaryExpression) Accept(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
}

func (binary *BinaryExpression) AcceptRecursive(visitor *Visitor) {
	visitor.VisitBinaryExpression(binary)
	binary.AcceptRecursive(visitor)
	binary.AcceptRecursive(visitor)
}

func (binary *BinaryExpression) Area() InputRegion {
	return binary.NodePosition
}


type CreateExpression struct {
	NodePosition InputRegion
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

func (create *CreateExpression) Area() InputRegion {
	return create.NodePosition
}


