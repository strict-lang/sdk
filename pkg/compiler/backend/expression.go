package backend

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (generation *Generation) GenerateIdentifier(identifier *tree.Identifier) {
	generation.Emit(identifier.Value)
}

func (generation *Generation) GenerateStringLiteral(literal *tree.StringLiteral) {
	generation.EmitFormatted(`"%s"`, literal.Value)
}

func (generation *Generation) GenerateNumberLiteral(literal *tree.NumberLiteral) {
	generation.Emit(literal.Value)
}

func (generation *Generation) GenerateBinaryExpression(binary *tree.BinaryExpression) {
	generation.EmitNode(binary.LeftOperand)
	generation.EmitFormatted(" %s ", binary.Operator.String())
	generation.EmitNode(binary.RightOperand)
}

func (generation *Generation) GenerateUnaryExpression(unary *tree.UnaryExpression) {
	generation.EmitFormatted("(%s", unary.Operator)
	generation.EmitNode(unary.Operand)
	generation.Emit(")")
}

func isPointerTarget(node tree.Node) bool {
	// TODO(merlinosayimwen): Replace this by attribute lookup
	if identifier, isIdentifier := node.(*tree.Identifier); isIdentifier {
		return identifier.Value == "this"
	}
	return false
}

func (generation *Generation) GenerateFieldSelectExpression(expression *tree.FieldSelectExpression) {
	if id, ok := expression.Target.(*tree.Identifier); ok {
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

func (generation *Generation) GenerateListSelectExpression(expression *tree.ListSelectExpression) {
	generation.EmitNode(expression.Target)
	generation.Emit("[")
	generation.EmitNode(expression.Index)
	generation.Emit("]")
}

func (generation *Generation) generateNamespaceSelector(selector *tree.FieldSelectExpression) {
	if generation.shouldInsertNamespaceSelector {
		generation.EmitNode(selector.Target)
		generation.Emit("::")
	}
	generation.EmitNode(selector.Selection)
}
