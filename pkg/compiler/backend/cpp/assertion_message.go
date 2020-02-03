package cpp

import (
	"strict.dev/sdk/pkg/compiler/grammar/token"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strings"
)

type AssertionMessageComputation struct {
	buffer  *strings.Builder
	visitor tree.Visitor
}

func NewAssertionMessageComputation() *AssertionMessageComputation {
	computation := &AssertionMessageComputation{}
	visitor := tree.NewEmptyVisitor()
	visitor.BinaryExpressionVisitor = computation.visitBinaryExpression
	visitor.UnaryExpressionVisitor = computation.visitUnaryExpression
	visitor.IdentifierVisitor = computation.visitIdentifier
	visitor.NumberLiteralVisitor = computation.visitNumberLiteral
	visitor.StringLiteralVisitor = computation.visitStringLiteral
	computation.visitor = visitor
	computation.buffer = &strings.Builder{}
	return computation
}

func ComputeAssertionMessage(assertedExpression tree.Node) string {
	computation := NewAssertionMessageComputation()
	computation.GenerateNode(assertedExpression)
	return computation.String()
}

func (computation *AssertionMessageComputation) String() string {
	return computation.buffer.String()
}

func (computation *AssertionMessageComputation) GenerateNode(node tree.Node) {
	node.Accept(computation.visitor)
}

func (computation *AssertionMessageComputation) visitBinaryExpression(expression *tree.BinaryExpression) {
	computation.buffer.WriteString("(")
	computation.GenerateNode(expression.LeftOperand)
	computation.buffer.WriteString(" ")
	message := translateComparisonOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ")
	computation.GenerateNode(expression.RightOperand)
	computation.buffer.WriteString(")")
}

func (computation *AssertionMessageComputation) visitUnaryExpression(expression *tree.UnaryExpression) {
	computation.buffer.WriteString("(")
	computation.GenerateNode(expression.Operand)
	computation.buffer.WriteString(" ")
	message := translateUnaryOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(")")
}

func (computation *AssertionMessageComputation) visitIdentifier(identifier *tree.Identifier) {
	computation.buffer.WriteString(identifier.Value)
}

func (computation *AssertionMessageComputation) visitStringLiteral(literal *tree.StringLiteral) {
	computation.buffer.WriteString(literal.Value)
}

func (computation *AssertionMessageComputation) visitNumberLiteral(literal *tree.NumberLiteral) {
	computation.buffer.WriteString(literal.Value)
}

var comparisonOperatorErrorMessages = map[token.Operator]string{
	token.EqualsOperator:        "is not equal to",
	token.NotEqualsOperator:     "is equal to",
	token.GreaterEqualsOperator: "is smaller than",
	token.GreaterOperator:       "is not greater than",
	token.SmallerEqualsOperator: "is greater than",
	token.SmallerOperator:       "is not smaller than",
}

func translateComparisonOperatorToErrorMessage(operator token.Operator) string {
	if message, found := comparisonOperatorErrorMessages[operator]; found {
		return message
	}
	return operator.String()
}

func translateUnaryOperatorToErrorMessage(operator token.Operator) string {
	if operator != token.NegateOperator {
		return operator.String()
	}
	return "is true"
}
