package codegen

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
)

func (generator *CodeGenerator) GenerateBinaryExpression(binary *ast.BinaryExpression) {
	binary.LeftOperand.Accept(generator.generators)
	generator.Emitf(" %s ", binary.Operator.String())
	binary.RightOperand.Accept(generator.generators)
}
