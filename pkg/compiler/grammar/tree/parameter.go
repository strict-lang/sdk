package tree

import "strict.dev/sdk/pkg/compiler/input"

type Parameter struct {
	Type   TypeName
	Name   *Identifier
	Region input.Region
	Parent Node
}

func (parameter *Parameter) SetEnclosingNode(target Node) {
  parameter.Parent = target
}

func (parameter *Parameter) EnclosingNode() (Node, bool) {
  return parameter.Parent, parameter.Parent != nil
}

func (parameter *Parameter) Accept(visitor Visitor) {
	visitor.VisitParameter(parameter)
}

func (parameter *Parameter) AcceptRecursive(visitor Visitor) {
	parameter.Accept(visitor)
	parameter.Name.AcceptRecursive(visitor)
	parameter.Type.AcceptRecursive(visitor)
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
