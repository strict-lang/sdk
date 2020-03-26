package cpp

import "gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"

type HeaderFileGeneration struct {
	Extension

	generation *Generation
}

func NewHeaderFileGeneration() *HeaderFileGeneration {
	return &HeaderFileGeneration{}
}

func (generation *HeaderFileGeneration) ModifyVisitor(
	parent *Generation, visitor *tree.DelegatingVisitor) {

	generation.generation = parent
	visitor.ClassDeclarationVisitor = generation.generateClassDeclaration
	generation.emitPragmas()
}

func (generation *HeaderFileGeneration) emitPragmas() {
	generation.generation.Emit("#pragma once")
	generation.generation.EmitEndOfLine()
	generation.generation.EmitEndOfLine()
}

func (generation *HeaderFileGeneration) generateClassDeclaration(
	declaration *tree.ClassDeclaration) {

	definition := newClassDefinition(generation.generation, declaration)
	definition.generateCode()
}
