package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type Parameter struct {
	Type   TypeName
	Name   *Identifier
	Region input.Region
}

func (parameter *Parameter) Accept(visitor Visitor) {
	VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor Visitor) {
	VisitParameter(parameter)
}

func (parameter *Parameter) Locate() input.Region {
	return parameter.Region
}
