package backend

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

func (generation *Generation) GenerateIdentifier(identifier *syntaxtree.Identifier) {
	generation.Emit(identifier.Value)
}

func (generation *Generation) GenerateStringLiteral(literal *syntaxtree.StringLiteral) {
	generation.EmitFormatted(`"%s"`, literal.Value)
}

func (generation *Generation) GenerateNumberLiteral(literal *syntaxtree.NumberLiteral) {
	generation.Emit(literal.Value)
}

func (generation *Generation) GenerateBinaryExpression(binary *syntaxtree.BinaryExpression) {
	generation.EmitNode(binary.LeftOperand)
	generation.EmitFormatted(" %s ", binary.Operator.String())
	generation.EmitNode(binary.RightOperand)
}

func (generation *Generation) GenerateUnaryExpression(unary *syntaxtree.UnaryExpression) {
	generation.EmitFormatted("(%s", unary.Operator)
	generation.EmitNode(unary.Operand)
	generation.Emit(")")
}

func isPointerTarget(node syntaxtree.Node) bool {
	// TODO(merlinosayimwen): Replace this by attribute lookup
	if identifier, isIdentifier := node.(*syntaxtree.Identifier); isIdentifier {
		return identifier.Value == "this"
	}
	return false
}

func (generation *Generation) GenerateSelectExpression(expression *syntaxtree.SelectExpression) {
	if id, ok := expression.Target.(*syntaxtree.Identifier); ok {
		if _, moduleExists := generation.importModules[id.Value]; moduleExists {
			generation.generateNamespaceSelector(expression)
			return
		}
	}
	generation.EmitNode(expression.Target)
	if isPointerTarget(expression.Target) {
		generation.Emit("->")
	} else {
		generation.Emit(".")
	}
	generation.EmitNode(expression.Selection)
}

func (generation *Generation) GenerateListSelectExpression(expression *syntaxtree.ListSelectExpression) {
	generation.EmitNode(expression.Target)
	generation.Emit("[")
	generation.EmitNode(expression.Index)
	generation.Emit("]")
}

func (generation *Generation) generateNamespaceSelector(selector *syntaxtree.SelectExpression) {
	if generation.shouldInsertNamespaceSelector {
		generation.EmitNode(selector.Target)
		generation.Emit("::")
	}
	generation.EmitNode(selector.Selection)
}
