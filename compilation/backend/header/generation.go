package header

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type Generation struct {
	backend.Extension

	generation *backend.Generation
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *ast.Visitor) {
	generation.generation = parent
	visitor.VisitClassDeclaration = generation.generateClassDeclaration
}

func (generation *Generation) generateClassDeclaration(declaration *ast.ClassDeclaration) {
	definition := newClassDefinition(generation.generation, declaration)
	definition.generateCode()
}
