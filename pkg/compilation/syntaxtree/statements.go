package syntaxtree

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) Position() Position {
	return expression.Expression.Position()
}

type BlockStatement struct {
	Children     []Node
	NodePosition Position
}

func (block *BlockStatement) Accept(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
}

func (block *BlockStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
	for _, statement := range block.Children {
		statement.AcceptRecursive(visitor)
	}
}

func (block *BlockStatement) Position() Position {
	return block.NodePosition
}

type ConditionalStatement struct {
	Condition    Node
	Alternative  Node
	Consequence  Node
	NodePosition Position
}

func (conditional *ConditionalStatement) HasAlternative() bool {
	return conditional.Alternative != nil
}

func (conditional *ConditionalStatement) Accept(visitor *Visitor) {
	visitor.VisitConditionalStatement(conditional)
}

func (conditional *ConditionalStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitConditionalStatement(conditional)
	conditional.Condition.AcceptRecursive(visitor)
	conditional.Consequence.AcceptRecursive(visitor)
	if conditional.HasAlternative() {
		conditional.Alternative.AcceptRecursive(visitor)
	}
}

func (conditional *ConditionalStatement) Position() Position {
	return conditional.NodePosition
}

// RangedLoopStatement is a control statement that. Counting from an initial
// value to some target while incrementing a field each step. The values of a
// ranged loop are numeral.
type RangedLoopStatement struct {
	NodePosition Position
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
	loop.InitialValue.AcceptRecursive(visitor)
	loop.EndValue.AcceptRecursive(visitor)
	loop.Body.AcceptRecursive(visitor)
}

func (loop *RangedLoopStatement) Position() Position {
	return loop.NodePosition
}

// ForEachLoopStatement is a control statement. Iterating an enumeration without
// requiring explicit indexing. As opposed to the ranged loop, the element
// iterated may be of any type that implements the 'Sequence' interface.
type ForEachLoopStatement struct {
	NodePosition Position
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
	loop.Body.AcceptRecursive(visitor)
}

func (loop *ForEachLoopStatement) Position() Position {
	return loop.NodePosition
}

type IncrementStatement struct {
	Operand      Node
	NodePosition Position
}

func (increment *IncrementStatement) Accept(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
}

func (increment *IncrementStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
	increment.Operand.AcceptRecursive(visitor)
}

func (increment *IncrementStatement) Position() Position {
	return increment.NodePosition
}

type DecrementStatement struct {
	Operand      Node
	NodePosition Position
}

func (decrement *DecrementStatement) Accept(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
}

func (decrement *DecrementStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
	decrement.Operand.AcceptRecursive(visitor)
}

func (decrement *DecrementStatement) Position() Position {
	return decrement.NodePosition
}

// YieldStatement yields an expression to an implicit list that is returned by
// the method it is defined in. Yield statements can only be in methods,
// returning a 'Sequence'. And their values type have to be of the sequences
// element type. Those statements are not accompanied by a ReturnStatement.
type YieldStatement struct {
	NodePosition Position
	Value        Node
}

func (yield *YieldStatement) Accept(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
}

func (yield *YieldStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
	yield.Value.AcceptRecursive(visitor)
}

func (yield *YieldStatement) Position() Position {
	return yield.NodePosition
}

// ReturnStatement is a control statement that can prematurely end the execution
// of a method or emit the return value. Return statements with a return value
// can only be defined in methods not returning 'void'. This statement is always
// the last statement in a block.
type ReturnStatement struct {
	NodePosition Position
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
	statement.Value.AcceptRecursive(visitor)
}

func (statement *ReturnStatement) Position() Position {
	return statement.NodePosition
}

// InvalidStatement represents a statement that has not been parsed correctly.
type InvalidStatement struct {
	NodePosition Position
}

func (statement *InvalidStatement) Accept(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) Position() Position {
	return statement.NodePosition
}

// EmptyStatement is a statement that does not execute any instructions.
type EmptyStatement struct {
	NodePosition Position
}

func (statement *EmptyStatement) Accept(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) Position() Position {
	return statement.NodePosition
}

// AssignStatement assigns values to left-hand-side expressions. Operations like
// add-assign are also represented by this Node. If the 'Target' node is a
// FieldDeclaration, this is a field definition.
type AssignStatement struct {
	Target       Node
	Value        Node
	Operator     token.Operator
	NodePosition Position
}

func (statement *AssignStatement) Accept(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
}

func (statement *AssignStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
	statement.Target.AcceptRecursive(visitor)
	statement.Value.AcceptRecursive(visitor)
}

func (statement *AssignStatement) Position() Position {
	return statement.NodePosition
}

type AssertStatement struct {
	NodePosition Position
	Expression   Node
}

func (assert *AssertStatement) Accept(visitor *Visitor) {
	visitor.VisitAssertStatement(assert)
}

func (assert *AssertStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitAssertStatement(assert)
	assert.Expression.AcceptRecursive(visitor)
}

func (assert *AssertStatement) Position() Position {
	return assert.NodePosition
}

type TestStatement struct {
	NodePosition Position
	Statements   Node
	MethodName   string
}

func (test *TestStatement) Accept(visitor *Visitor) {
	visitor.VisitTestStatement(test)
}

func (test *TestStatement) AcceptRecursive(visitor *Visitor) {
	visitor.VisitTestStatement(test)
	test.Statements.AcceptRecursive(visitor)
}

func (test *TestStatement) Position() Position {
	return test.NodePosition
}
