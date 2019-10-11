package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type BlockStatement struct {
	Children     []Statement
	Region input.Region
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
