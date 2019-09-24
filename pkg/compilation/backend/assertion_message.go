package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"strings"
)

type AssertionMessageComputation struct {
	buffer  strings.Builder
	visitor *syntaxtree2.Visitor
}

func NewAssertionMessageComputation() *AssertionMessageComputation {
	computation := &AssertionMessageComputation{}
	visitor := syntaxtree2.NewEmptyVisitor()
	visitor.VisitBinaryExpression = computation.visitBinaryExpression
	visitor.VisitUnaryExpression = computation.visitUnaryExpression
	visitor.VisitIdentifier = computation.visitIdentifier
	visitor.VisitNumberLiteral = computation.visitNumberLiteral
	visitor.VisitStringLiteral = computation.visitStringLiteral
	computation.visitor = visitor
	return computation
}

func ComputeAssertionMessage(assertedExpression syntaxtree2.Node) string {
	computation := NewAssertionMessageComputation()
	computation.GenerateNode(assertedExpression)
	return computation.String()
}

func (computation *AssertionMessageComputation) String() string {
	return computation.buffer.String()
}

func (computation *AssertionMessageComputation) GenerateNode(node syntaxtree2.Node) {
	node.Accept(computation.visitor)
}

func (computation *AssertionMessageComputation) visitBinaryExpression(expression *syntaxtree2.BinaryExpression) {
	computation.buffer.WriteString("( ")
	computation.GenerateNode(expression.LeftOperand)
	computation.buffer.WriteString(" ")
	message := translateComparisonOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ")
	computation.GenerateNode(expression.RightOperand)
	computation.buffer.WriteString(") ")
}

func (computation *AssertionMessageComputation) visitUnaryExpression(expression *syntaxtree2.UnaryExpression) {
	computation.buffer.WriteString("( ")
	computation.GenerateNode(expression.Operand)
	computation.buffer.WriteString(" ")
	message := translateUnaryOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ) ")
}

func (computation *AssertionMessageComputation) visitIdentifier(identifier *syntaxtree2.Identifier) {
	computation.buffer.WriteString(identifier.Value)
}

func (computation *AssertionMessageComputation) visitStringLiteral(literal *syntaxtree2.StringLiteral) {
	computation.buffer.WriteString(literal.Value)
}

func (computation *AssertionMessageComputation) visitNumberLiteral(literal *syntaxtree2.NumberLiteral) {
	computation.buffer.WriteString(literal.Value)
}

var comparisonOperatorErrorMessages = map[token2.Operator]string{
	token2.EqualsOperator:        "is not equal to",
	token2.NotEqualsOperator:     "is equal to",
	token2.GreaterEqualsOperator: "is smaller than",
	token2.GreaterOperator:       "is not greater than",
	token2.SmallerEqualsOperator: "is greater than",
	token2.SmallerOperator:       "is not smaller than",
}

func translateComparisonOperatorToErrorMessage(operator token2.Operator) string {
	if message, found := comparisonOperatorErrorMessages[operator]; found {
		return message
	}
	return operator.String()
}

func translateUnaryOperatorToErrorMessage(operator token2.Operator) string {
	if operator != token2.NegateOperator {
		return operator.String()
	}
	return "is true"
}
