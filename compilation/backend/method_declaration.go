package backend

import "gitlab.com/strict-lang/sdk/compilation/ast"

func (generation *Generation) EmitMethodDeclaration(declaration *ast.MethodDeclaration) {
	returnTypeName := updateTypeName(declaration.Type)
	methodName := declaration.Name.Value
	generation.EmitIndent()
	generation.EmitFormatted("%s %s(", returnTypeName.FullName(), methodName)
	for index, parameter := range declaration.Parameters {
		if index != 0 {
			generation.Emit(", ")
		}
		parameterTypeName := updateTypeName(parameter.Type)
		generation.EmitFormatted("%s %s", parameterTypeName.FullName(), parameter.Name.Value)
	}
	generation.Emit(")")
}
