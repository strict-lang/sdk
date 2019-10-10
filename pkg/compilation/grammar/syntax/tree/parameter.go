package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

type Parameter struct {
	Type         TypeName
	Name         *Identifier
	Region  input.Region
}

func (parameter *Parameter) Accept(visitor Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) Locate() input.Region {
	return parameter.Region
}
