package headerfile

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Generation struct {
	backend.Extension

	generation *backend.Generation
}

func NewGeneration() *Generation {
	return &Generation{}
}

func (generation *Generation) ModifyVisitor(parent *backend.Generation, visitor *tree.DelegatingVisitor) {
	generation.generation = parent
	visitor.ClassDeclarationVisitor = generation.generateClassDeclaration
	generation.emitPragmas()
}

func (generation *Generation) emitPragmas() {
	generation.generation.Emit("#pragma once")
	generation.generation.EmitEndOfLine()
	generation.generation.EmitEndOfLine()
}

func (generation *Generation) generateClassDeclaration(declaration *tree.ClassDeclaration) {
	definition := newClassDefinition(generation.generation, declaration)
	definition.generateCode()
}
