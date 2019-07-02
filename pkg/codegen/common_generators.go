package codegen

import "github.com/BenjaminNitschke/Strict/pkg/ast"

func (generator *CodeGenerator) GenerateIdentifier(identifier *ast.Identifier) {
	generator.Emit(identifier.Value)
}

func (generator *CodeGenerator) GenerateStringLiteral(literal *ast.StringLiteral) {
	generator.Emit(literal.Value)
}

func (generator *CodeGenerator) GenerateNumberLiteral(literal *ast.NumberLiteral) {
	generator.Emit(literal.Value)
}