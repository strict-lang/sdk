package codegen

import "gitlab.com/strict-lang/sdk/compiler/ast"

func (generator *CodeGenerator) GenerateIdentifier(identifier *ast.Identifier) {
	generator.Emit(identifier.Value)
}

func (generator *CodeGenerator) GenerateSelectExpression(access *ast.SelectorExpression) {
	generator.EmitNode(access.Target)
	generator.Emit(".")
	generator.EmitNode(access.Selection)
}

func (generator *CodeGenerator) GenerateStringLiteral(literal *ast.StringLiteral) {
	generator.Emitf(`"%s"`, literal.Value)
}

func (generator *CodeGenerator) GenerateNumberLiteral(literal *ast.NumberLiteral) {
	generator.Emit(literal.Value)
}

func (generator *CodeGenerator) GenerateExpressionStatement(statement *ast.ExpressionStatement) {
	generator.EmitNode(statement.Expression)
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
		generator.EmitNode(child)
	}
	generator.leaveBlock()
	generator.Emit("\n")
	generator.Spaces()
	generator.Emit("}")
}
