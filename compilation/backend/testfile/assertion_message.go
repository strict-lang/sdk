package testfile

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/token"
	"strings"
)

type assertionMessageComputation struct {
	buffer strings.Builder
	visitor *ast.Visitor
}

func newAssertionMessageComputation() *assertionMessageComputation {
	computation := &assertionMessageComputation{}
	visitor := ast.NewEmptyVisitor()
	visitor.VisitBinaryExpression = computation.visitBinaryExpression
	visitor.VisitUnaryExpression = computation.visitUnaryExpression
	visitor.VisitIdentifier = computation.visitIdentifier
	visitor.VisitNumberLiteral = computation.visitNumberLiteral
	visitor.VisitStringLiteral = computation.visitStringLiteral
	computation.visitor = visitor
	return computation
}

func (computation *assertionMessageComputation) String() string {
	return computation.buffer.String()
}

func (computation *assertionMessageComputation) generateNode(node ast.Node) {
	node.Accept(computation.visitor)
}

func (computation *assertionMessageComputation) visitBinaryExpression(expression *ast.BinaryExpression) {
	computation.buffer.WriteString("( ")
	computation.generateNode(expression.LeftOperand)
	computation.buffer.WriteString(" ")
	message := translateComparisonOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ")
	computation.generateNode(expression.RightOperand)
	computation.buffer.WriteString(") ")
}

func (computation *assertionMessageComputation) visitUnaryExpression(expression *ast.UnaryExpression) {
	computation.buffer.WriteString("( ")
	computation.generateNode(expression.Operand)
	computation.buffer.WriteString(" ")
	message := translateUnaryOperatorToErrorMessage(expression.Operator)
	computation.buffer.WriteString(message)
	computation.buffer.WriteString(" ) ")
}

func (computation *assertionMessageComputation) visitIdentifier(identifier *ast.Identifier) {
	computation.buffer.WriteString(identifier.Value)
}

func (computation *assertionMessageComputation) visitStringLiteral(literal *ast.StringLiteral) {
	computation.buffer.WriteString(literal.Value)
}

func (computation *assertionMessageComputation) visitNumberLiteral(literal *ast.NumberLiteral) {
	computation.buffer.WriteString(literal.Value)
}

var comparisonOperatorErrorMessages = map[token.Operator] string {
	token.EqualsOperator: "is not equal to",
	token.NotEqualsOperator: "is equal to",
	token.GreaterEqualsOperator: "is smaller than",
	token.GreaterOperator: "is not greater than",
	token.SmallerEqualsOperator: "is greater than",
	token.SmallerOperator: "is not smaller than",
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