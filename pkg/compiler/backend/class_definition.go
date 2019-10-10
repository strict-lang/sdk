package backend

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/grammar/syntax/tree"
)

func (generation *Generation) GenerateClassDeclaration(declaration *tree.ClassDeclaration) {
	for _, child := range declaration.Children {
		generation.EmitNode(child)
	}
}
