package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
)

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Expression
	RightOperand Expression
	Operator     token.Operator
	Region       input.Region
	resolvedType resolvedType
	Parent       Node
}

func (binary *BinaryExpression) SetEnclosingNode(target Node) {
	binary.Parent = target
}

func (binary *BinaryExpression) EnclosingNode() (Node, bool) {
	return binary.Parent, binary.Parent != nil
}

func (binary *BinaryExpression) ResolveType(class *scope.Class) {
	binary.resolvedType.resolve(class)
}

func (binary *BinaryExpression) ResolvedType() (*scope.Class, bool) {
	return binary.resolvedType.class()
}

func (binary *BinaryExpression) Accept(visitor Visitor) {
	visitor.VisitBinaryExpression(binary)
}

func (binary *BinaryExpression) AcceptRecursive(visitor Visitor) {
	binary.Accept(visitor)
	binary.LeftOperand.AcceptRecursive(visitor)
	binary.RightOperand.AcceptRecursive(visitor)
}

func (binary *BinaryExpression) Locate() input.Region {
	return binary.Region
}

func (binary *BinaryExpression) Matches(node Node) bool {
	if target, ok := node.(*BinaryExpression); ok {
		return binary.matchesExpression(target)
	}
	return false
}

func (binary *BinaryExpression) matchesExpression(target *BinaryExpression) bool {
	return binary.Operator == target.Operator &&
		binary.LeftOperand.Matches(target.LeftOperand) &&
		binary.RightOperand.Matches(target.RightOperand)
}

func (binary *BinaryExpression) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteBinaryExpression(binary)
}

func (binary *BinaryExpression) TransformExpressions(transformer ExpressionTransformer) {
	binary.LeftOperand = binary.LeftOperand.Transform(transformer)
	binary.RightOperand = binary.RightOperand.Transform(transformer)
}
