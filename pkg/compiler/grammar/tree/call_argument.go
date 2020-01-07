package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type CallArgumentList []*CallArgument

type CallArgument struct {
	Label  string
	Value  Node
	Region input.Region
	Parent Node
}

func (argument *CallArgument) SetEnclosingNode(target Node) {
  argument.Parent = target
}

func (argument *CallArgument) EnclosingNode() (Node, bool) {
  return argument.Parent, argument.Parent != nil
}

func (argument *CallArgument) IsLabeled() bool {
	return argument.Label != ""
}

func (argument *CallArgument) Accept(visitor Visitor) {
	visitor.VisitCallArgument(argument)
}

func (argument *CallArgument) AcceptRecursive(visitor Visitor) {
	argument.Accept(visitor)
	argument.Value.AcceptRecursive(visitor)
}

func (argument *CallArgument) Locate() input.Region {
	return argument.Region
}

func (argument *CallArgument) Matches(node Node) bool {
	if target, ok := node.(*CallArgument); ok {
		return argument.Label == target.Label &&
			argument.Value.Matches(target.Value)
	}
	return false
}
