package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type Parameter struct {
	Type   TypeName
	Name   *Identifier
	Region input.Region
}

func (parameter *Parameter) Accept(visitor Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor Visitor) {
	parameter.Accept(visitor)
}

func (parameter *Parameter) Locate() input.Region {
	return parameter.Region
}

func (parameter *Parameter) Matches(node Node) bool {
	if target, ok := node.(*Parameter); ok {
		return parameter.Name.Matches(target.Name) &&
			parameter.Type.Matches(target.Type)
	}
	return false
}