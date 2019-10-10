package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
)

type Statement interface {
	Node
	IsModifyingControlFlow() bool
}

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) Area() InputRegion {
	return InputRegion()
}

type BlockStatement struct {
	Children     []Node
	NodePosition InputRegion
}

func (block *BlockStatement) Accept(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
}

func (block *BlockStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
	for _, statement := range block.Children {
		AcceptRecursive(visitor)
	}
}

func (block *BlockStatement) Area() InputRegion {
	return block.NodePosition
}

// RangedLoopStatement is a control statement that. Counting from an initial
// value to some target while incrementing a field each step. The values of a
// ranged loop are numeral.
type RangedLoopStatement struct {
	NodePosition InputRegion
	ValueField   *Identifier
	InitialValue Node
	EndValue     Node
	Body         Node
}

func (loop *RangedLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitRangedLoopStatement(loop)
}

func (loop *RangedLoopStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitRangedLoopStatement(loop)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (loop *RangedLoopStatement) Area() InputRegion {
	return loop.NodePosition
}

// ForEachLoopStatement is a control statement. Iterating an enumeration without
// requiring explicit indexing. As opposed to the ranged loop, the element
// iterated may be of any type that implements the 'Sequence' interface.
type ForEachLoopStatement struct {
	NodePosition InputRegion
	Body         Node
	Field        *Identifier
	Sequence     Node
}

func (loop *ForEachLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitForEachLoopStatement(loop)
}

func (loop *ForEachLoopStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitForEachLoopStatement(loop)
	loop.Field.AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (loop *ForEachLoopStatement) Area() InputRegion {
	return loop.NodePosition
}

type IncrementStatement struct {
	Operand      Node
	NodePosition InputRegion
}

func (increment *IncrementStatement) Accept(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
}

func (increment *IncrementStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
	AcceptRecursive(visitor)
}

func (increment *IncrementStatement) Area() InputRegion {
	return increment.NodePosition
}

type DecrementStatement struct {
	Operand      Node
	NodePosition InputRegion
}

func (decrement *DecrementStatement) Accept(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
}

func (decrement *DecrementStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
	AcceptRecursive(visitor)
}

func (decrement *DecrementStatement) Area() InputRegion {
	return decrement.NodePosition
}

// YieldStatement yields an expression to an implicit list that is returned by
// the method it is defined in. Yield statements can only be in methods,
// returning a 'Sequence'. And their values type have to be of the sequences
// element type. Those statements are not accompanied by a ReturnStatement.
type YieldStatement struct {
	NodePosition InputRegion
	Value        Node
}

func (yield *YieldStatement) Accept(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
}

func (yield *YieldStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
	AcceptRecursive(visitor)
}

func (yield *YieldStatement) Area() InputRegion {
	return yield.NodePosition
}

// ReturnStatement is a control statement that can prematurely end the execution
// of a method or emit the return value. Return statements with a return value
// can only be defined in methods not returning 'void'. This statement is always
// the last statement in a block.
type ReturnStatement struct {
	NodePosition InputRegion
	Value        Node
}

func (statement *ReturnStatement) IsReturningValue() bool {
	return statement.Value != nil
}

func (statement *ReturnStatement) Accept(visitor *Visitor) {
	visitor.VisitReturnStatement(statement)
}

func (statement *ReturnStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitReturnStatement(statement)
	AcceptRecursive(visitor)
}

func (statement *ReturnStatement) Area() InputRegion {
	return statement.NodePosition
}

// InvalidStatement represents a statement that has not been parsed correctly.
type InvalidStatement struct {
	NodePosition InputRegion
}

func (statement *InvalidStatement) Accept(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) Area() InputRegion {
	return statement.NodePosition
}

// EmptyStatement is a statement that does not execute any instructions.
type EmptyStatement struct {
	NodePosition InputRegion
}

func (statement *EmptyStatement) Accept(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) Area() InputRegion {
	return statement.NodePosition
}

// AssignStatement assigns values to left-hand-side expressions. Operations like
// add-assign are also represented by this Node. If the 'Target' node is a
// FieldDeclaration, this is a field definition.
type AssignStatement struct {
	Target       Node
	Value        Node
	Operator     token.Operator
	NodePosition InputRegion
}

func (statement *AssignStatement) Accept(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
}

func (statement *AssignStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
	AcceptRecursive(visitor)
	AcceptRecursive(visitor)
}

func (statement *AssignStatement) Area() InputRegion {
	return statement.NodePosition
}
