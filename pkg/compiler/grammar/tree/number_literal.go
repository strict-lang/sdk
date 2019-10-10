package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type NumberLiteral struct {
	Value  string
	Region input.Region
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
