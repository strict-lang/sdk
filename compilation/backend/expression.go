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

func (generation *Generation) GenerateSelectExpression(expression *ast.SelectExpression) {
	if id, ok := expression.Target.(*ast.Identifier); ok {
		if _, moduleExists := generation.importModules[id.Value]; moduleExists {
			generation.generateNamespaceSelector(expression)
			return
		}
	}
	generation.EmitNode(expression.Target)
	generation.Emit(".")
	generation.EmitNode(expression.Selection)
}

func (generation *Generation) GenerateListSelectExpression(expression *ast.ListSelectExpression) {
	generation.EmitNode(expression.Target)
	generation.Emit("[")
	generation.EmitNode(expression.Index)
	generation.Emit("]")
}

func (generation *Generation) generateNamespaceSelector(selector *ast.SelectExpression) {
	generation.EmitNode(selector.Target)
	generation.Emit("::")
	generation.EmitNode(selector.Selection)
}
