package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

func (generator *CodeGenerator) GenerateIdentifier(identifier *ast.Identifier) {
	generator.Emit(identifier.Value)
}

func (generator *CodeGenerator) GenerateStringLiteral(literal *ast.StringLiteral) {
	generator.Emitf(`"%s"`, literal.Value)
}

func (generator *CodeGenerator) GenerateNumberLiteral(literal *ast.NumberLiteral) {
	generator.Emit(literal.Value)
}

func (generator *CodeGenerator) GenerateExpressionStatement(statement *ast.ExpressionStatement) {
	statement.Expression.Accept(generator.generators)
	generator.Emit(";")
}

func (generator *CodeGenerator) GenerateBlockStatement(block *ast.BlockStatement) {
	generator.Emit("{\n")
	generator.enterBlock()
	for index, child := range block.Children {
		if index != 0 {
			generator.Emit("\n")
		}
		generator.Spaces()
		child.Accept(generator.generators)
	}
	generator.leaveBlock()
	generator.Emit("\n")
	generator.Spaces()
	generator.Emit("}")
}