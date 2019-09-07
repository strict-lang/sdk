package backend

import "gitlab.com/strict-lang/sdk/compilation/ast"

func (generation *Generation) EmitMethodDeclaration(declaration *ast.MethodDeclaration) {
	generation.EmitIndent()

	generation.EmitNode(declaration.Type)
	generation.Emit(" ")
	generation.EmitNode(declaration.Name)

	generation.Emit("(")
	for index, parameter := range declaration.Parameters {
		if index != 0 {
			generation.Emit(", ")
		}
		generation.EmitNode(parameter)
	}
	generation.Emit(")")
}

func (generation *Generation) GenerateParameter(parameter *ast.Parameter) {
	generation.EmitNode(parameter.Type)
	generation.Emit(" ")
	generation.EmitNode(parameter.Name)
}

