package cpp

import (
	"strict.dev/sdk/pkg/compiler/grammar/tree"
)

func (generation *Generation) GenerateClassDeclaration(declaration *tree.ClassDeclaration) {
	for _, child := range declaration.Children {
		generation.EmitNode(child)
	}
}
