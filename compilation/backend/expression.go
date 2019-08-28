package backend

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
)

func (generation *Generation) GenerateIdentifier(identifier *ast.Identifier) {
	generation.Emit(identifier.Value)
}

func (generation *Generation) GenerateStringLiteral(literal *ast.StringLiteral) {
	generation.EmitFormatted(`"%s"`, literal.Value)
}

func (generation *Generation) GenerateNumberLiteral(literal *ast.NumberLiteral) {
	generation.Emit(literal.Value)
}

func (generation *Generation) GenerateBinaryExpression(binary *ast.BinaryExpression) {
	generation.EmitNode(binary.LeftOperand)
	generation.EmitFormatted(" %s ", binary.Operator.String())
	generation.EmitNode(binary.RightOperand)
}

func (generation *Generation) GenerateUnaryExpression(unary *ast.UnaryExpression) {
	generation.EmitFormatted("(%s", unary.Operator)
	generation.EmitNode(unary.Operand)
	generation.Emit(")")
}

func (generation *Generation) GenerateSelectorExpression(selector *ast.SelectorExpression) {
	if id, ok := selector.Target.(*ast.Identifier); ok {
		if _, moduleExists := generation.importModules[id.Value]; moduleExists {
			generation.generateNamespaceSelector(selector)
			return
		}
	}
	generation.EmitNode(selector.Target)
	generation.Emit(".")
	generation.EmitNode(selector.Selection)
}

func (generation *Generation) generateNamespaceSelector(selector *ast.SelectorExpression) {
	generation.EmitNode(selector.Target)
	generation.Emit("::")
	generation.EmitNode(selector.Selection)
}
