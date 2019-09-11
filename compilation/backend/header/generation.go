package header

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type Generation struct{
	generation *backend.Generation
}

func (generation *Generation)

func (generation *Generation) generateClassDeclaration(declaration *ast.ClassDeclaration) {
	definition := newClassDefinition(generation.generation, declaration)
	definition.generateCode()
}