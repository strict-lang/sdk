package backend

import "gitlab.com/strict-lang/sdk/compilation/ast"

func (generation *Generation) GenerateClassDeclaration(declaration *ast.ClassDeclaration) {
	for _, child := range declaration.Children {
		generation.EmitNode(child)
	}
}