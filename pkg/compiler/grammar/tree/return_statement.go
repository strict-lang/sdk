package tree

import "github.com/strict-lang/sdk/pkg/compiler/input"

// ReturnStatement is a control statement that can prematurely end the execution
// of a method or emit the return value. Return statements with a return value
// can only be defined in methods not returning 'void'. This statement is always
// the last statement in a block.
type ReturnStatement struct {
	Region input.Region
	Value  Expression
	Parent Node
}

func (statement *ReturnStatement) SetEnclosingNode(target Node) {
	statement.Parent = target
}

func (statement *ReturnStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}

func (statement *ReturnStatement) IsReturningValue() bool {
	return statement.Value != nil
}

func (statement *ReturnStatement) Accept(visitor Visitor) {
	visitor.VisitReturnStatement(statement)
}

func (statement *ReturnStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	if statement.IsReturningValue() {
		statement.Value.AcceptRecursive(visitor)
	}
}

func (statement *ReturnStatement) Locate() input.Region {
	return statement.Region
}

func (statement *ReturnStatement) Matches(node Node) bool {
	if target, ok := node.(*ReturnStatement); ok {
		return statement.matchesReturnValue(target)
	}
	return false
}

func (statement *ReturnStatement) matchesReturnValue(
	target *ReturnStatement) bool {

	if statement.Value == nil {
		return target.Value == nil
	}
	return statement.Value.Matches(target.Value)
}

func (statement *ReturnStatement) TransformExpressions(transformer ExpressionTransformer) {
	if statement.Value != nil {
		statement.Value = statement.Value.Transform(transformer)
	}
}
