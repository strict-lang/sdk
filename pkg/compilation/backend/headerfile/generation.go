package headerfile

import (
	backend2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Generation struct {
	backend2.Extension

	generation *backend2.Generation
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend2.Generation, visitor *syntaxtree2.Visitor) {
	generation.generation = parent
	visitor.VisitClassDeclaration = generation.generateClassDeclaration
	generation.emitPragmas()
}

func (generation *Generation) emitPragmas() {
	generation.generation.Emit("#pragma once")
	generation.generation.EmitEndOfLine()
	generation.generation.EmitEndOfLine()
}

func (generation *Generation) generateClassDeclaration(declaration *syntaxtree2.ClassDeclaration) {
	definition := newClassDefinition(generation.generation, declaration)
	definition.generateCode()
}
