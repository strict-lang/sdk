package cpp

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
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

func (generation *Generation) GenerateFieldSelectExpression(expression *tree.ChainExpression) {
	if id, ok := expression.FirstChild().(*tree.Identifier); ok {
		if _, moduleExists := generation.importModules[id.Value]; moduleExists {
			generation.generateNamespaceSelector(expression)
			return
		}
	}
	remaining := expression.Expressions[1:]
	lastIndex := len(remaining) - 1
	for index, remaining := range remaining {
		generation.EmitNode(remaining)
		if index != lastIndex {
			if isPointerTarget(remaining) {
				generation.Emit("->")
			} else {
				generation.Emit(".")
			}
		}
	}
}

func (generation *Generation) GenerateListSelectExpression(expression *tree.ListSelectExpression) {
	generation.EmitNode(expression.Target)
	generation.Emit("[")
	generation.EmitNode(expression.Index)
	generation.Emit("]")
}

func (generation *Generation) generateNamespaceSelector(selector *tree.ChainExpression) {
	if generation.shouldInsertNamespaceSelector {
		generation.EmitNode(selector.FirstChild())
		generation.Emit("::")
	}
	generation.EmitNode(selector.LastChild())
}
