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

func (block *StatementBlock) ReplaceExact(replaced Node, target Statement) {
	for index, child := range block.Children {
		if child == replaced {
			block.Children[index] = target
		}
	}
}

func (block *StatementBlock) ReplaceMatching(filter NodeFilter, target Statement) {
	for index, child := range block.Children {
		if filter(child) {
			block.Children[index] = target
		}
	}
}

func (block *StatementBlock) FindIndexOfExact(node Statement) (int, bool) {
	for index, child := range block.Children {
		if child == node {
			return index, true
		}
	}
	return 0, false
}

func (block *StatementBlock) FindIndexOfMatching(filter NodeFilter) (int, bool) {
	for index, child := range block.Children {
		if filter(child) {
			return index, true
		}
	}
	return 0, false
}

func (block *StatementBlock) InsertBeforeOffset(offset input.Offset, node Statement) {
	for index, child := range block.Children {
		if child.Locate().Begin() < offset {
			block.Children[index] = node
		}
	}
}

func (block *StatementBlock) InsertBeforeIndex(index int, node Statement) {
	if index <= 1 {
		block.Prepend(node)
	} else {
		insert(block.Children, index-1, node)
	}
}

func (block *StatementBlock) InsertAfterIndex(index int, node Statement) {
	if index >= len(block.Children) {
		block.Append(node)
	} else {
		insert(block.Children, index, node)
	}
}

func insert(slice []Statement, index int, value Statement) (result []Statement) {
	result = append(slice, nil)
	copy(result[index + 1:], result[index:])
	result[index] = value
	return result
}

func (block *StatementBlock) Prepend(node Statement) {
	newHead := []Statement{node}
	block.Children = append(newHead, block.Children...)
}

func (block *StatementBlock) Append(node Statement) {
	block.Children = append(block.Children, node)
}