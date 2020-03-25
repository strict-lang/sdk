package cpp

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

const InitMethodName = "Generated$Init"

func (generation *Generation) EmitMethodDeclaration(declaration *tree.MethodDeclaration) {
	generation.EmitIndent()
	generation.EmitNode(declaration.Type)
	generation.Emit(" ")
	generation.EmitNode(declaration.Name)
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
