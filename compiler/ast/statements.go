package ast

import (
	"gitlab.com/strict-lang/sdk/compiler/token"
	"strings"
)

type ExpressionStatement struct {
	Expression Node
}

func (expression *ExpressionStatement) Accept(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitExpressionStatement(expression)
}

func (expression *ExpressionStatement) Position() Position {
	return expression.Expression.Position()
}

type MethodCall struct {
	// Method is the called method. It can be any kind of expression
	// with the value of a method. Common nodes are identifiers and
	// field selectors.
	Method    Node
	// An array of expression nodes that are the arguments passed to
	// the method. The arguments types are checked during type checking.
	Arguments []Node
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

type BlockStatement struct {
	Children []Node
	NodePosition Position
}

func (block *BlockStatement) Accept(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
}

func (block *BlockStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitBlockStatement(block)
	for _, statement := range block.Children {
		statement.AcceptAll(visitor)
	}
}

func (block *BlockStatement) Position() Position {
	return block.NodePosition
}

type ConditionalStatement struct {
	Condition Node
	Alternative Node
	Consequence Node
	NodePosition Position
}

func (conditional *ConditionalStatement) HasAlternative() bool {
	return conditional.Alternative != nil
}

func (conditional *ConditionalStatement) Accept(visitor *Visitor) {
	visitor.VisitConditionalStatement(conditional)
}

func (conditional *ConditionalStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitConditionalStatement(conditional)
	conditional.Condition.AcceptAll(visitor)
	conditional.Consequence.AcceptAll(visitor)
	if conditional.HasAlternative() {
		conditional.Alternative.AcceptAll(visitor)
	}
}

func (conditional *ConditionalStatement) Position() Position {
	return conditional.NodePosition
}

// Loop statement that counts from an initial value to some target
// while incrementing the current value each step. The values of a
// ranged loop are numeral.
type RangedLoopStatement struct {
	NodePosition Position
	// Name of the field in which the current value is stored.
	ValueField *Identifier
	// Initial value assigned to the value field.
	InitialValue Node
	// Value that, when reached, breaks the loop.
	EndValue Node
	// Body of the loop that is executed each time after the value field
	// is updated. May contain break and continue statements.
	Body  Node
}

func (loop *RangedLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitRangedLoopStatement(loop)
}

func (loop *RangedLoopStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitRangedLoopStatement(loop)
	loop.InitialValue.AcceptAll(visitor)
	loop.EndValue.AcceptAll(visitor)
	loop.Body.AcceptAll(visitor)
}

func (loop *RangedLoopStatement) Position() Position {
	return loop.NodePosition
}

// Loop that iterates an enumeration of elements without requiring
// explicit indexing. As opposed to the ranged loop, the element
// iterated may be of any type.
type ForEachLoopStatement struct {
	NodePosition Position
	// Body of the loop that is executed for every element in the collection.
	// May contain break and continue statements.
	Body Node
	// Field is the name of the local field that has the value of
	// the current element of target.
	Field *Identifier
	// Target is the collection that is iterated.
	Enumeration Node
}

func (loop *ForEachLoopStatement) Accept(visitor *Visitor) {
	visitor.VisitForEachLoopStatement(loop)
}

func (loop *ForEachLoopStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitForEachLoopStatement(loop)
	loop.Enumeration.AcceptAll(visitor)
	loop.Body.AcceptAll(visitor)
}

func (loop *ForEachLoopStatement) Position() Position {
	return loop.NodePosition
}

type IncrementStatement struct {
	Operand Node
	NodePosition Position
}

func (increment *IncrementStatement) Accept(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
}

func (increment *IncrementStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitIncrementStatement(increment)
	increment.Operand.AcceptAll(visitor)
}

func (increment *IncrementStatement) Position() Position {
	return increment.NodePosition
}

type DecrementStatement struct {
	Operand Node
	NodePosition Position
}

func (decrement *DecrementStatement) Accept(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
}

func (decrement *DecrementStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitDecrementStatement(decrement)
	decrement.Operand.AcceptAll(visitor)
}

func (decrement *DecrementStatement) Position() Position {
	return decrement.NodePosition
}

type YieldStatement struct {
	NodePosition Position
	// Value is the value that is yielded.
	Value Node
}

func (yield *YieldStatement) Accept(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
}

func (yield *YieldStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitYieldStatement(yield)
	yield.Value.AcceptAll(visitor)
}

func (yield *YieldStatement) Position() Position {
	return yield.NodePosition
}

type ReturnStatement struct {
	NodePosition Position
	// Value is the value that is returned.
	// This pointer can be nil, if the return does not return a value.
	Value Node
}

func (statement *ReturnStatement) Accept(visitor *Visitor) {
	visitor.VisitReturnStatement(statement)
}

func (statement *ReturnStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitReturnStatement(statement)
	statement.Value.AcceptAll(visitor)
}

func (statement *ReturnStatement) Position() Position {
	return statement.NodePosition
}

type InvalidStatement struct {
	NodePosition Position
}

func (statement *InvalidStatement) Accept(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitInvalidStatement(statement)
}

func (statement *InvalidStatement) Position() Position {
	return statement.NodePosition
}

type EmptyStatement struct {
	NodePosition Position
}

func (statement *EmptyStatement) Accept(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitEmptyStatement(statement)
}

func (statement *EmptyStatement) Position() Position {
	return statement.NodePosition
}

type AssignStatement struct {
	Target   Node
	Value    Node
	Operator token.Operator
	NodePosition Position
}

func (statement *AssignStatement) Accept(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
}

func (statement *AssignStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitAssignStatement(statement)
	statement.Target.AcceptAll(visitor)
	statement.Value.AcceptAll(visitor)
}

func (statement *AssignStatement) Position() Position {
	return statement.NodePosition
}

type ImportStatement struct {
	Path string
	Alias *Identifier
	NodePosition Position
}

func (statement *ImportStatement) ModuleName() string {
	if statement.Alias == nil || statement.Alias.Value == "" {
		return moduleNameByPath(statement.Path)
	}
	return statement.Alias.Value
}

func moduleNameByPath(path string) string {
	var begin = 0
	if strings.Contains(path, "/") {
		begin = strings.LastIndex(path, "/") + 1
	}
	var end int
	if strings.HasSuffix(path, ".h") {
		end = len(path) - 2
	} else {
		end = len(path)
	}
	return path[begin:end]
}
func (statement *ImportStatement) Accept(visitor *Visitor) {
	visitor.VisitImportStatement(statement)
}

func (statement *ImportStatement) AcceptAll(visitor *Visitor) {
	visitor.VisitImportStatement(statement)
}

func (statement *ImportStatement) Position() Position {
	return statement.NodePosition
}
