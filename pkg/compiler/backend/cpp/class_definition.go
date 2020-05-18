package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (generation *Generation) GenerateClassDeclaration(declaration *tree.ClassDeclaration) {
	generation.isGeneratingApp = isApp(declaration)
	if generation.isGeneratingApp {
		generation.runMethod = findRunMethod(declaration)
	}
	for _, child := range declaration.Children {
		generation.EmitNode(child)
	}
}
