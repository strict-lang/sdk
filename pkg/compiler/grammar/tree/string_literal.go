package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"strconv"
)

type StringLiteral struct {
	// Value is the strings value. It does not contain a leading and trailing
	// character. Escaped strings are C like.
	Value string
	// Region is the region of the input that contain the literal.
	// It contains the leading and trailing characters.
	Region input.Region
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
		Region: literal.NodeRegion,
	}, nil
}
