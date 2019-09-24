package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

func (generation *Generation) GenerateClassDeclaration(declaration *syntaxtree2.ClassDeclaration) {
	for _, child := range declaration.Children {
		generation.EmitNode(child)
	}
}
