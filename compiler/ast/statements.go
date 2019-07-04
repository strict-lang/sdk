package ast

import (
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

type MethodCall struct {
	// Method is the called method. It can be any kind of expression
	// with the value of a method. Common nodes are identifiers and
	// field selectors.
	Method    Node
	Arguments []Node
}

func (call *MethodCall) Accept(visitor *Visitor) {
	visitor.VisitMethodCall(call)
}

type BlockStatement struct {
	Children []Node
}

func (block *BlockStatement) Accept(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
}

type ConditionalStatement struct {
	Else      Node
	Body      Node
	Condition Node
}

func (conditional *ConditionalStatement) Accept(visitor *Visitor) {
	visitor.VisitConditionalStatement(conditional)
}

type FromToLoopStatement struct {
	Body  Node
	Field Identifier
	From  Node
	To    Node
}

func (loop *FromToLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitFromToLoopStatement(loop)
}

type ForeachLoopStatement struct {
	Body Node
	// Field is the name of the local field that has the value of
	// the current element of target.
	Field Identifier
	// Target is the collection that is iterated.
	Target Node
}

func (loop *ForeachLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitForeachLoopStatement(loop)
}

type IncrementStatement struct {
	Operand Node
}

type DecrementStatement struct {
	Operand Node
}

func (increment *IncrementStatement) Accept(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
}

func (decrement *DecrementStatement) Accept(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
}

type YieldStatement struct {
	// Value is the value that is yielded.
	Value Node
}

func (yield *YieldStatement) Accept(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
}

type ReturnStatement struct {
	// Value is the value that is returned.
	// This pointer can be nil, if the return does not return a value.
	Value Node
}

func (statement *ReturnStatement) Accept(visitor *Visitor) {
	visitor.VisitReturnStatement(statement)
}

type InvalidStatement struct{}

func (statement *InvalidStatement) Accept(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

type EmptyStatement struct{}

func (statement *EmptyStatement) Accept(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

type AssignStatement struct {
	Target   Node
	Value    Node
	Operator token.Operator
}

func (statement *AssignStatement) Accept(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
}
