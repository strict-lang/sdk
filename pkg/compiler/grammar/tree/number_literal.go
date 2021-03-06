package tree

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"strconv"
	"strings"
)

type NumberLiteral struct {
	Value        string
	Region       input.Region
	Parent       Node
	resolvedType resolvedType
}

const floatBitSize = 64

func (literal *NumberLiteral) SetEnclosingNode(target Node) {
	literal.Parent = target
}

func (literal *NumberLiteral) EnclosingNode() (Node, bool) {
	return literal.Parent, literal.Parent != nil
}

func (literal *NumberLiteral) ResolveType(class *scope.Class) {
	literal.resolvedType.resolve(class)
}

func (literal *NumberLiteral) ResolvedType() (*scope.Class, bool) {
	return literal.resolvedType.class()
}

func (literal *NumberLiteral) IsFloat() bool {
	_, err := strconv.ParseFloat(literal.Value, floatBitSize)
	return err == nil && strings.Contains(literal.Value, ".")
}

const integerBitSize = 32
const decimal = 10

func (literal *NumberLiteral) AsInt() (int, bool) {
	number, err := strconv.ParseInt(literal.Value, decimal, integerBitSize)
	return int(number), err == nil
}

func (literal *NumberLiteral) Accept(visitor Visitor) {
	visitor.VisitNumberLiteral(literal)
}

func (literal *NumberLiteral) AcceptRecursive(visitor Visitor) {
	literal.Accept(visitor)
}

func (literal *NumberLiteral) Locate() input.Region {
	return literal.Region
}

func (literal *NumberLiteral) Matches(node Node) bool {
	if target, ok := node.(*NumberLiteral); ok {
		return literal.Value == target.Value
	}
	return false
}

func (literal *NumberLiteral) Transform(transformer ExpressionTransformer) Expression {
	return transformer.RewriteNumberLiteral(literal)
}
