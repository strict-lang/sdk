package ast

import (
	"fmt"
	"strings"
)

type Printing struct {
	buffer strings.Builder
	indent int
	visitor *Visitor
}

func Print(node Node) {
	printing := &Printing{}
	visitor := &Visitor{
		VisitParameter:            nil,
		VisitMethodCall:           nil,
		VisitIdentifier:           printing.printIdentifier,
		VisitTestStatement:        nil,
		VisitStringLiteral:        printing.printStringLiteral,
		VisitNumberLiteral:        printing.printNumberLiteral,
		VisitEmptyStatement:       nil,
		VisitYieldStatement:       nil,
		VisitBlockStatement:       nil,
		VisitAssertStatement:      nil,
		VisitUnaryExpression:      printing.printUnaryExpression,
		VisitImportStatement:      nil,
		VisitAssignStatement:      nil,
		VisitReturnStatement:      nil,
		VisitTranslationUnit:      printing.printTranslationUnit,
		VisitCreateExpression:     nil,
		VisitInvalidStatement:     nil,
		VisitFieldDeclaration:     nil,
		VisitGenericTypeName:      nil,
		VisitConcreteTypeName:     nil,
		VisitBinaryExpression:     printing.printBinaryExpression,
		VisitMethodDeclaration:    nil,
		VisitSelectorExpression:   nil,
		VisitIncrementStatement:   nil,
		VisitDecrementStatement:   nil,
		VisitRangedLoopStatement:  nil,
		VisitExpressionStatement:  printing.printExpressionStatement,
		VisitForEachLoopStatement: nil,
		VisitConditionalStatement: nil,
	}
	printing.visitor = visitor
	printing.printNode(node)
	fmt.Println(printing.buffer.String())
}

func (printing *Printing) print(message string) {
	printing.buffer.WriteString(message)
}

func (printing *Printing) printLine(message string) {
	printing.print(message)
	printing.printNewLine()
}

func (printing *Printing) increaseIndent() {
	printing.indent++
}

func (printing *Printing) decreaseIndent() {
	printing.indent--
}

func (printing *Printing) printFieldName(name string) {
	printing.printFormatted("%s: ", name)
}

func (printing *Printing) printIndentedFieldName(name string) {
	printing.printIndent()
	printing.printFieldName(name)
}

func (printing *Printing) printIndentedNodeField(name string, node Node) {
	printing.printIndentedFieldName(name)
	printing.printNode(node)
	printing.printNewLine()
}

func (printing *Printing) printIndentedStringField(name string, value string) {
	printing.printIndentedFieldName(name)
	printing.print(value)
	printing.printNewLine()
}

func (printing *Printing) printFormatted(message string, arguments ...interface{}) {
	printing.buffer.WriteString(fmt.Sprintf(message, arguments...))
}

func (printing *Printing) printIndent() {
	for count := 0; count < printing.indent; count++ {
		printing.print("\t")
	}
}

func (printing *Printing) printNewLine() {
	printing.print("\n")
}

func (printing *Printing) printNode(node Node) {
	node.Accept(printing.visitor)
}

func (printing *Printing) printBinaryExpression(expression *BinaryExpression) {
	printing.printLine("BinaryExpression:")
	printing.increaseIndent()
	printing.printIndentedStringField("operator",expression.Operator.String())
	printing.printIndentedNodeField("leftOperand", expression.LeftOperand)
	printing.printIndentedNodeField("rightOperand", expression.RightOperand)
	printing.decreaseIndent()
}

func (printing *Printing) printUnaryExpression(expression *UnaryExpression) {
	printing.printLine("UnaryExpression:")
	printing.increaseIndent()
	printing.printIndentedStringField("operator", expression.Operator.String())
	printing.printIndentedNodeField("operand", expression.Operand)
	printing.decreaseIndent()
}

func (printing *Printing) printExpressionStatement(statement *ExpressionStatement) {
	printing.print("!")
	printing.printNode(statement.Expression)
}

func (printing *Printing) printTranslationUnit(unit *TranslationUnit) {
	printing.printLine("TranslationUnit:")
	printing.increaseIndent()
	printing.printIndentedStringField("name", unit.name)
	for _, node := range unit.Children {
		printing.printIndent()
		printing.print("- ")
		printing.printNode(node)
	}
	printing.decreaseIndent()
}

func (printing *Printing) printIdentifier(identifier *Identifier) {
	printing.printFormatted("%s", identifier.Value)
}

func (printing *Printing) printStringLiteral(literal *StringLiteral) {
	printing.printFormatted("\"%s\"", literal.Value)
}

func (printing *Printing) printNumberLiteral(number *NumberLiteral) {
	printing.print(number.Value)
}