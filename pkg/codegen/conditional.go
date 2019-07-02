package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

func (generator *CodeGenerator) GenerateConditionalStatement(statement *ast.ConditionalStatement) {
	generator.Emit("if (")
	statement.Condition.Accept(generator.generators)
	generator.Emit(") {")
	statement.Body.Accept(generator.generators)
	generator.Emit("}")
	if statement.Else != nil {
		generator.Emit("else ")
		statement.Else.Accept(generator.generators)
	}
}
