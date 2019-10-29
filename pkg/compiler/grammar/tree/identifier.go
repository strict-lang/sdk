package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type Identifier struct {
	Value  string
	Region input.Region
	resolvedType resolvedType
}


func (identifier *Identifier) Accept(visitor Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) AcceptRecursive(visitor Visitor) {
	identifier.Accept(visitor)
}

func (identifier *Identifier) Locate() input.Region {
	return identifier.Region
}
