package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

// BinaryExpression is an operation on two operands.
type BinaryExpression struct {
	LeftOperand  Node
	RightOperand Node
	Operator     token.Operator
	Region       input.Region
	resolvedType resolvedType
}

func (binary *BinaryExpression) GetResolvedType() (TypeDescriptor, bool) {
	return binary.resolvedType.getDescriptor()
}

func (binary *BinaryExpression) Resolve(descriptor TypeDescriptor) {
	binary.resolvedType.setDescriptor(descriptor)
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
