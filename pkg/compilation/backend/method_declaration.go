package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

const InitMethodName = "Generated$Init"

func (generation *Generation) EmitMethodDeclaration(declaration *syntaxtree2.MethodDeclaration) {
	generation.EmitIndent()
	generation.EmitNode(declaration.Type)
	generation.Emit(" ")
	generation.EmitNode(declaration.Name)
	generation.EmitParameterList(declaration.Parameters)
}

func (generation *Generation) EmitParameterList(parameters syntaxtree2.ParameterList) {
	generation.Emit("(")
	for index, parameter := range parameters {
		if index != 0 {
			generation.Emit(", ")
		}
		generation.EmitNode(parameter)
	}
	generation.Emit(")")
}

func (generation *Generation) GenerateParameter(parameter *syntaxtree2.Parameter) {
	generation.EmitNode(parameter.Type)
	generation.Emit(" ")
	generation.EmitNode(parameter.Name)
}
