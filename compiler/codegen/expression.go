package codegen

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
)

func (generator *CodeGenerator) GenerateBinaryExpression(binary *ast.BinaryExpression) {
	binary.LeftOperand.Accept(generator.generators)
	generator.Emitf(" %s ", binary.Operator.String())
	binary.RightOperand.Accept(generator.generators)
}

func (generator *CodeGenerator) GenerateUnaryExpression(unary *ast.UnaryExpression) {
	generator.Emitf("(%s", unary.Operator)
	unary.Operand.Accept(generator.generators)
	generator.Emit(")")
}