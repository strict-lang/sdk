package backend

import "gitlab.com/strict-lang/sdk/compilation/syntaxtree"

func (generation *Generation) GenerateClassDeclaration(declaration *syntaxtree.ClassDeclaration) {
	for _, child := range declaration.Children {
		generation.EmitNode(child)
	}
}
