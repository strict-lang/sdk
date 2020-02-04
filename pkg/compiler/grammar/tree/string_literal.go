package tree

import (
	"strconv"
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/scope"
)

type StringLiteral struct {
	// Value is the strings value. It does not contain a leading and trailing
	// character. Escaped strings are C like.
	Value string
	// Region is the region of the input that contain the literal.
	// It contains the leading and trailing characters.
	Region       input.Region
	resolvedType resolvedType
	Parent       Node
}

func (literal *StringLiteral) SetEnclosingNode(target Node) {
	literal.Parent = target
}

func (literal *StringLiteral) EnclosingNode() (Node, bool) {
	return literal.Parent, literal.Parent != nil
}

func (literal *StringLiteral) ResolveType(class *scope.Class) {
	literal.resolvedType.resolve(class)
}

func (literal *StringLiteral) ResolvedType() (*scope.Class, bool) {
	return literal.resolvedType.class()
}

func (literal *StringLiteral) Accept(visitor Visitor) {
	visitor.VisitStringLiteral(literal)
}

func (literal *StringLiteral) AcceptRecursive(visitor Visitor) {
	literal.Accept(visitor)
}

func (literal *StringLiteral) Locate() input.Region {
	return literal.Region
}

func (literal *StringLiteral) ToStringLiteral() (*StringLiteral, error) {
	return literal, nil
}

const floatBitSizeLimit = 64

func (literal *StringLiteral) ToNumberLiteral() (*NumberLiteral, error) {
	_, err := strconv.ParseFloat(literal.Value, floatBitSizeLimit)
	if err != nil {
		return nil, err
	}
	return &NumberLiteral{
		Value:  literal.Value,
		Region: literal.Region,
	}, nil
}

func (literal *StringLiteral) Matches(node Node) bool {
	if target, ok := node.(*StringLiteral); ok {
		return literal.Value == target.Value
	}
	return false
}

func (literal *StringLiteral) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteStringLiteral(literal)
}
