package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

type StatementBlock struct {
	Children []Statement
	Region   input.Region
	Parent   Node
	scope    scope.Scope
}

func (block *StatementBlock) UpdateScope(target scope.Scope) {
  block.scope = target
}

func (block *StatementBlock) Scope() scope.Scope {
  return block.scope
}

func (block *StatementBlock) SetEnclosingNode(target Node) {
	block.Parent = target
}

func (block *StatementBlock) EnclosingNode() (Node, bool) {
  return block.Parent, block.Parent != nil
}

func (block *StatementBlock) Accept(visitor Visitor) {
	visitor.VisitBlockStatement(block)
}

func (block *StatementBlock) AcceptRecursive(visitor Visitor) {
	block.Accept(visitor)
	for _, statement := range block.Children {
		statement.AcceptRecursive(visitor)
	}
}

func (block *StatementBlock) Locate() input.Region {
	return block.Region
}

func (block *StatementBlock) Matches(node Node) bool {
	if target, ok := node.(*StatementBlock); ok {
		return block.hasChildren(target.Children)
	}
	return false
}

func (block *StatementBlock) hasChildren(children []Statement) bool {
	if len(block.Children) != len(children) {
		return false
	}
	for index := 0; index < len(children); index++ {
		if !block.Children[index].Matches(children[index]) {
			return false
		}
	}
	return true
}
