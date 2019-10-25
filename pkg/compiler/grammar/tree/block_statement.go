package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type BlockStatement struct {
	Children []Statement
	Region   input.Region
}

func (block *BlockStatement) Accept(visitor Visitor) {
	visitor.VisitBlockStatement(block)
}

func (block *BlockStatement) AcceptRecursive(visitor Visitor) {
	block.Accept(visitor)
	for _, statement := range block.Children {
		statement.AcceptRecursive(visitor)
	}
}

func (block *BlockStatement) Locate() input.Region {
	return block.Region
}

func (block *BlockStatement) Matches(node Node) bool {
	if target, ok := node.(*BlockStatement); ok {
		return block.hasChildren(target.Children)
	}
	return false
}

func (block *BlockStatement) hasChildren(children []Statement) bool {
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