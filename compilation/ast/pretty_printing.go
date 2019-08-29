package ast

import (
	"fmt"
	"strings"
)

type Printing struct {
	buffer  strings.Builder
	indent  int
	visitor *Visitor
}

func Print(node Node) {
	printing := &Printing{}
	visitor := &Visitor{
		VisitParameter:            printing.printParameter,
		VisitMethodCall:           printing.printMethodCall,
		VisitIdentifier:           printing.printIdentifier,
		VisitTestStatement:        printing.printTestStatement,
		VisitStringLiteral:        printing.printStringLiteral,
		VisitNumberLiteral:        printing.printNumberLiteral,
		VisitEmptyStatement:       printing.printEmptyStatement,
		VisitYieldStatement:       printing.printYieldStatement,
		VisitBlockStatement:       printing.printBlockStatement,
		VisitAssertStatement:      printing.printAssertStatement,
		VisitUnaryExpression:      printing.printUnaryExpression,
		VisitImportStatement:      printing.printImportStatement,
		VisitAssignStatement:      printing.printAssignStatement,
		VisitReturnStatement:      printing.printReturnStatement,
		VisitTranslationUnit:      printing.printTranslationUnit,
		VisitCreateExpression:     printing.printCreateExpression,
		VisitInvalidStatement:     printing.printInvalidStatement,
		VisitFieldDeclaration:     printing.printFieldDeclaration,
		VisitGenericTypeName:      printing.printGenericTypeName,
		VisitConcreteTypeName:     printing.printConcreteTypeName,
		VisitBinaryExpression:     printing.printBinaryExpression,
		VisitMethodDeclaration:    printing.printMethodDeclaration,
		VisitSelectorExpression:   printing.printSelectorExpression,
		VisitIncrementStatement:   printing.printIncrementStatement,
		VisitDecrementStatement:   printing.printDecrementStatement,
		VisitRangedLoopStatement:  printing.printRangedLoopStatement,
		VisitExpressionStatement:  printing.printExpressionStatement,
		VisitForEachLoopStatement: printing.printForEachLoopStatement,
		VisitConditionalStatement: printing.printConditionalStatement,
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

func (printing *Printing) printNodeBegin(name string) {
	printing.printLine(name + ": ")
	printing.increaseIndent()
}

func (printing *Printing) printNodeEnd() {
	printing.decreaseIndent()
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
	printing.printNodeBegin("BinaryExpression")
	printing.printIndentedStringField("operator", expression.Operator.String())
	printing.printIndentedNodeField("leftOperand", expression.LeftOperand)
	printing.printIndentedNodeField("rightOperand", expression.RightOperand)
	printing.printNodeEnd()
}

func (printing *Printing) printUnaryExpression(expression *UnaryExpression) {
	printing.printNodeBegin("UnaryExpression")
	printing.printIndentedStringField("operator", expression.Operator.String())
	printing.printIndentedNodeField("operand", expression.Operand)
	printing.printNodeEnd()
}

func (printing *Printing) printExpressionStatement(statement *ExpressionStatement) {
	printing.print("!")
	printing.printNode(statement.Expression)
}

func (printing *Printing) printTranslationUnit(unit *TranslationUnit) {
	printing.printNodeBegin("TranslationUnit")
	printing.printIndentedStringField("name", unit.name)
	for _, node := range unit.Children {
		printing.printIndent()
		printing.print("- ")
		printing.printNode(node)
	}
	printing.printNodeEnd()
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

func (printing *Printing) printAssertStatement(statement *AssertStatement) {
	printing.printNodeBegin("AssertStatement")
	printing.printIndentedNodeField("expression", statement.Expression)
	printing.printNodeEnd()
}

func (printing *Printing) printReturnStatement(statement *ReturnStatement) {
	printing.printNodeBegin("ReturnStatement")
	if statement.Value == nil {
		printing.printIndentedStringField("value", "None")
	} else {
		printing.printIndentedNodeField("value", statement.Value)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printYieldStatement(statement *YieldStatement) {
	printing.printNodeBegin("YieldStatement")
	printing.printIndentedNodeField("value", statement.Value)
	printing.printNodeEnd()
}

func (printing *Printing) printConditionalStatement(statement *ConditionalStatement) {
	printing.printNodeBegin("ConditionalStatement")
	printing.printIndentedNodeField("condition", statement.Condition)
	printing.printIndentedNodeField("consequence", statement.Consequence)
	if statement.HasAlternative() {
		printing.printIndentedNodeField("Alternative", statement.Alternative)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printRangedLoopStatement(statement *RangedLoopStatement) {
	printing.printNodeBegin("RangedLoopStatement")
	printing.printIndentedNodeField("valueField", statement.ValueField)
	printing.printIndentedNodeField("initialValue", statement.InitialValue)
	printing.printIndentedNodeField("endValue", statement.EndValue)
	printing.printIndentedNodeField("body", statement.Body)
	printing.printNodeEnd()
}

func (printing *Printing) printForEachLoopStatement(statement *ForEachLoopStatement) {
	printing.printNodeBegin("ForEachLoopStatement")
	printing.printIndentedNodeField("field", statement.Field)
	printing.printIndentedNodeField("enumeration", statement.Enumeration)
	printing.printIndentedNodeField("body", statement.Body)
	printing.printNodeEnd()
}

func (printing *Printing) printMethodDeclaration(method *MethodDeclaration) {
	printing.printNodeBegin("MethodDeclaration")
	printing.printIndentedNodeField("name", method.Name)
	printing.printIndentedNodeField("returnType", method.Type)
	printing.printFieldName("parameters")
	printing.printNewLine()
	for _, parameter := range method.Parameters {
		printing.printIndent()
		printing.print("- ")
		printing.printNode(parameter)
		printing.printNewLine()
	}
	if method.Body != nil {
		printing.printIndentedNodeField("body", method.Body)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printParameter(parameter *Parameter) {
	printing.printNodeBegin("Parameter")
	printing.printIndentedNodeField("name", parameter.Name)
	printing.printIndentedNodeField("type", parameter.Type)
	printing.printNodeEnd()
}

func (printing *Printing) printGenericTypeName(name *GenericTypeName) {
	printing.print(name.FullName())
}

func (printing *Printing) printConcreteTypeName(name *ConcreteTypeName) {
	printing.print(name.FullName())
}

func (printing *Printing) printInvalidStatement(statement *InvalidStatement) {
	printing.print("!!!INVALID")
}

func (printing *Printing) printEmptyStatement(statement *EmptyStatement) {
	printing.print("-")
}

func (printing *Printing) printIncrementStatement(statement *IncrementStatement) {
	printing.printNodeBegin("IncrementStatement")
	printing.printIndentedNodeField("operand", statement.Operand)
	printing.printNodeEnd()
}

func (printing *Printing) printDecrementStatement(statement *DecrementStatement) {
	printing.printNodeBegin("DecrementStatement")
	printing.printIndentedNodeField("operand", statement.Operand)
	printing.printNodeEnd()
}

func (printing *Printing) printSelectorExpression(expression *SelectorExpression) {
	printing.printNodeBegin("Selector")
	printing.printIndentedNodeField("target", expression.Target)
	printing.printIndentedNodeField("selection", expression.Selection)
	printing.printNodeEnd()
}

func (printing *Printing) printMethodCall(call *MethodCall) {
	printing.printNodeBegin("MethodCall")
	printing.printIndentedNodeField("method", call.Method)
	printing.printFieldName("parameters")
	printing.printNewLine()
	for _, argument := range call.Arguments {
		printing.printIndent()
		printing.print("- ")
		printing.printNode(argument)
		printing.printNewLine()
	}
	printing.printNodeEnd()
}

func (printing *Printing) printFieldDeclaration(field *FieldDeclaration) {
	printing.printNodeBegin("FieldDeclaration")
	printing.printIndentedNodeField("name", field.Name)
	printing.printIndentedNodeField("type", field.TypeName)
	printing.printNodeEnd()
}

func (printing *Printing) printAssignStatement(statement *AssignStatement) {
	printing.printNodeBegin("AssignStatement")
	printing.printIndentedNodeField("target", statement.Target)
	printing.printIndentedNodeField("value", statement.Value)
	printing.printNodeEnd()
}

func (printing *Printing) printTestStatement(statement *TestStatement) {
	printing.printNodeBegin("TestStatement")
	printing.printIndentedStringField("methodName", statement.MethodName)
	printing.printIndentedNodeField("body", statement.Statements)
	printing.printNodeEnd()
}

func (printing *Printing) printBlockStatement(statement *BlockStatement) {
	printing.printNodeBegin("BlockStatement")
	for _, child := range statement.Children {
		printing.printIndent()
		printing.print("- ")
		printing.printNode(child)
		printing.printNewLine()
	}
	printing.printNodeEnd()
}

func (printing *Printing) printImportStatement(statement *ImportStatement) {
	printing.printNodeBegin("ImportStatement")
	printing.printIndentedStringField("path", statement.Path)
	if statement.Alias != nil {
		printing.printIndentedNodeField("alias", statement.Alias)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printCreateExpression(expression *CreateExpression) {
	printing.printNodeBegin("CreateExpression")
	printing.printIndentedNodeField("constructor", expression.Constructor)
	printing.printNodeEnd()
}