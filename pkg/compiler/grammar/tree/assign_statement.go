package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/input"
)

// AssignStatement assigns values to left-hand-side expressions. Operations like
// add-assign are also represented by this Node. If the 'Target' node is a
// FieldDeclaration, this is a field definition.
type AssignStatement struct {
	Target   Node
	Value    Expression
	Operator token.Operator
	Region   input.Region
	Parent   Node
}

func (assign *AssignStatement) SetEnclosingNode(target Node) {
	assign.Parent = target
}

func (assign *AssignStatement) EnclosingNode() (Node, bool) {
	return assign.Parent, assign.Parent != nil
}

func (assign *AssignStatement) Accept(visitor Visitor) {
	visitor.VisitAssignStatement(assign)
}

func (assign *AssignStatement) AcceptRecursive(visitor Visitor) {
	assign.Accept(visitor)
	assign.Target.AcceptRecursive(visitor)
	assign.Value.AcceptRecursive(visitor)
}

func (assign *AssignStatement) Locate() input.Region {
	return assign.Region
}

func (assign *AssignStatement) Matches(node Node) bool {
	if target, ok := node.(*AssignStatement); ok {
		return assign.matchesAssign(target)
	}
	return false
}

func (assign *AssignStatement) matchesAssign(target *AssignStatement) bool {
	return assign.Operator == target.Operator &&
		assign.Target.Matches(target.Target) &&
		assign.Value.Matches(target.Value)
}

func (assign *AssignStatement) TransformExpressions(transformer ExpressionTransformer) {
	assign.Value = assign.Value.Transform(transformer)
}
