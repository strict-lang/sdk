package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

const InitMethodName = "Generated$Init"

func (generation *Generation) EmitMethodDeclaration(declaration *tree.MethodDeclaration) {
	generation.EmitIndent()
	generation.emitPossiblyAbstractMethodSignature(declaration)
}

func (generation *Generation) emitPossiblyAbstractMethodSignature(
	declaration *tree.MethodDeclaration) {

	if generation.Unit.Class.Trait {
		generation.Emit("virtual ")
	}
	generation.emitMethodSignature(declaration)
}

func (generation *Generation) emitMethodSignature(declaration *tree.MethodDeclaration) {
	if declaration.Factory {
		generation.emitFactorySignature(declaration)
		return
	}
	generation.EmitNode(declaration.Type)
	generation.Emit(" ")
	generation.EmitNode(declaration.Name)
	generation.EmitParameterList(declaration.Parameters)
}

func (generation *Generation) emitFactorySignature(declaration *tree.MethodDeclaration) {
	generation.Emit("explicit ")
	generation.Emit(generation.Unit.Class.Name)
	generation.EmitParameterList(declaration.Parameters)
}

func (generation *Generation) EmitParameterList(parameters tree.ParameterList) {
	generation.Emit("(")
	for index, parameter := range parameters {
		if index != 0 {
			generation.Emit(", ")
		}
		generation.EmitNode(parameter)
	}
	generation.Emit(")")
}

func (generation *Generation) GenerateParameter(parameter *tree.Parameter) {
	generation.EmitNode(parameter.Type)
	generation.Emit(" ")
	generation.EmitNode(parameter.Name)
}
