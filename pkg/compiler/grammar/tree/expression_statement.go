package tree

import "github.com/strict-lang/sdk/pkg/compiler/input"

type ExpressionStatement struct {
	Expression Expression
	Parent     Node
}

func (statement *ExpressionStatement) SetEnclosingNode(target Node) {
	statement.Parent = target
}

func (statement *ExpressionStatement) EnclosingNode() (Node, bool) {
	return statement.Parent, statement.Parent != nil
}

func (statement *ExpressionStatement) Accept(visitor Visitor) {
	visitor.VisitExpressionStatement(statement)
}

func (statement *ExpressionStatement) AcceptRecursive(visitor Visitor) {
	statement.Accept(visitor)
	statement.Expression.AcceptRecursive(visitor)
}

func (statement *ExpressionStatement) Locate() input.Region {
	return statement.Expression.Locate()
}

func (statement *ExpressionStatement) Matches(node Node) bool {
	if target, ok := node.(*ExpressionStatement); ok {
		return statement.Expression.Matches(target.Expression)
	}
	return false
}

func (statement *ExpressionStatement) TransformExpressions(transformer ExpressionTransformer) {
	statement.Expression = statement.Expression.Transform(transformer)
}
