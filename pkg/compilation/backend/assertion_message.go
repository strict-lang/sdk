package backend

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"strings"
)

type AssertionMessageComputation struct {
	buffer  strings.Builder
	visitor *syntaxtree.Visitor
}

func NewAssertionMessageComputation() *AssertionMessageComputation {
	computation := &AssertionMessageComputation{}
	visitor := syntaxtree.NewEmptyVisitor()
	visitor.VisitBinaryExpression = computation.visitBinaryExpression
	visitor.VisitUnaryExpression = computation.visitUnaryExpression
	visitor.VisitIdentifier = computation.visitIdentifier
	visitor.VisitNumberLiteral = computation.visitNumberLiteral
	visitor.VisitStringLiteral = computation.visitStringLiteral
	computation.visitor = visitor
	return computation
}

func ComputeAssertionMessage(assertedExpression syntaxtree.Node) string {
	computation := NewAssertionMessageComputation()
	computation.GenerateNode(assertedExpression)
	return computation.String()
}

func (computation *AssertionMessageComputation) String() string {
	return computation.buffer.String()
}

func (computation *AssertionMessageComputation) GenerateNode(node syntaxtree.Node) {
	node.Accept(computation.visitor)
}

func (computation *AssertionMessageComputation) visitBinaryExpression(expression *syntaxtree.BinaryExpression) {
	computation.buffer.WriteString("( ")
	computation.GenerateNode(expression.LeftOperand)
	computation.buffer.WriteString(" ")
	message := translateComparisonOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ")
	computation.GenerateNode(expression.RightOperand)
	computation.buffer.WriteString(") ")
}

func (computation *AssertionMessageComputation) visitUnaryExpression(expression *syntaxtree.UnaryExpression) {
	computation.buffer.WriteString("( ")
	computation.GenerateNode(expression.Operand)
	computation.buffer.WriteString(" ")
	message := translateUnaryOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ) ")
}

func (computation *AssertionMessageComputation) visitIdentifier(identifier *syntaxtree.Identifier) {
	computation.buffer.WriteString(identifier.Value)
}

func (computation *AssertionMessageComputation) visitStringLiteral(literal *syntaxtree.StringLiteral) {
	computation.buffer.WriteString(literal.Value)
}

func (computation *AssertionMessageComputation) visitNumberLiteral(literal *syntaxtree.NumberLiteral) {
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
