package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
)

func (generator *CodeGenerator) GenerateBinaryExpression(binary *ast.BinaryExpression) {
	generator.EmitNode(binary.LeftOperand)
	generator.Emitf(" %s ", binary.Operator.String())
	generator.EmitNode(binary.RightOperand)
}

func (generator *CodeGenerator) GenerateUnaryExpression(unary *ast.UnaryExpression) {
	generator.Emitf("(%s", unary.Operator)
	generator.EmitNode(unary.Operand)
	generator.Emit(")")
}

func (generator *CodeGenerator) GenerateSelectorExpression(selector *ast.SelectorExpression) {
	if id, ok := selector.Target.(*ast.Identifier); ok {
		if _, moduleExists := generator.importModules[id.Value]; moduleExists {
			generator.generateNamespaceSelector(selector)
			return
		}
	}
	generator.EmitNode(selector.Target)
	generator.Emit(".")
	generator.EmitNode(selector.Selection)
}

func (generator *CodeGenerator) generateNamespaceSelector(selector *ast.SelectorExpression) {
	generator.EmitNode(selector.Target)
	generator.Emit("::")
	generator.EmitNode(selector.Selection)
}
