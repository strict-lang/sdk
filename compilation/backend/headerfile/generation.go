package headerfile

import (
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type Generation struct {
	backend.Extension

	generation *backend.Generation
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *syntaxtree.Visitor) {
	generation.generation = parent
	visitor.VisitClassDeclaration = generation.generateClassDeclaration
	visitor.VisitConstructorDeclaration = generation.generateConstructorDeclaration
}

func (generation *Generation) generateClassDeclaration(declaration *syntaxtree.ClassDeclaration) {
	definition := newClassDefinition(generation.generation, declaration)
	definition.generateCode()
}

func (generation *Generation) generateConstructorDeclaration(declaration *syntaxtree.ConstructorDeclaration) {
	output := generation.generation
	className := generation.generation.Unit.Class.Name
	output.Emit(className)
	output.EmitParameterList(declaration.Parameters)
}