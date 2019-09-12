package ast

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

const prettyPrintIndent = "  "

type Printing struct {
	buffer             strings.Builder
	indent             int
	visitor            *Visitor
	colored            bool
	noNewLineAfterNode int
}

func newPrinting() *Printing {
	printing := &Printing{}
	visitor := &Visitor{
		VisitParameter:            printing.printParameter,
		VisitMethodCall:           printing.printMethodCall,
		VisitIdentifier:           printing.printIdentifier,
		VisitListTypeName:         printing.printListTypeName,
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
		VisitClassDeclaration:     printing.printClassDeclaration,
		VisitGenericTypeName:      printing.printGenericTypeName,
		VisitConcreteTypeName:     printing.printConcreteTypeName,
		VisitBinaryExpression:     printing.printBinaryExpression,
		VisitMethodDeclaration:    printing.printMethodDeclaration,
		VisitSelectorExpression:   printing.printSelectExpression,
		VisitIncrementStatement:   printing.printIncrementStatement,
		VisitDecrementStatement:   printing.printDecrementStatement,
		VisitRangedLoopStatement:  printing.printRangedLoopStatement,
		VisitExpressionStatement:  printing.printExpressionStatement,
		VisitForEachLoopStatement: printing.printForEachLoopStatement,
		VisitConditionalStatement: printing.printConditionalStatement,
		VisitListSelectExpression: printing.printListSelectExpression,
	}
	printing.visitor = visitor
	return printing
}

func newColoredPrinting() *Printing {
	printing := newPrinting()
	printing.colored = true
	return printing
}

// Print prints a pretty representation of the AST node.
func Print(node Node) {
	printing := newPrinting()
	printing.printNode(node)
	fmt.Println(printing.buffer.String())
}

// PrintColored prints a pretty representation of the AST node, that uses
// colors to highlight the output in a way that makes it easier to read.
func PrintColored(node Node) {
	printing := newColoredPrinting()
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
	if printing.colored {
		coloredName := color.CyanString("\"%s\"", name)
		coloredAssign := color.RedString(" = ")
		printing.print(coloredName + coloredAssign)
		return
	}
	printing.printFormatted("\"%s\" = ", name)
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

func (printing *Printing) printIndentedListFieldBegin(name string) {
	printing.printIndentedFieldName(name)
	if printing.colored {
		printing.printLine(color.RedString("["))
	} else {
		printing.printLine("[")
	}
	printing.increaseIndent()
	printing.noNewLineAfterNode++
}

func (printing *Printing) printListFieldEnd() {
	printing.decreaseIndent()
	printing.printIndent()
	if printing.colored {
		printing.printLine(color.RedString("]"))
	} else {
		printing.printLine("]")
	}
}

func (printing *Printing) printNodeBegin(name string) {
	if printing.colored {
		coloredName := color.GreenString(name)
		coloredCurly := color.RedString(" {")
		printing.printLine(coloredName + coloredCurly)
	} else {
		printing.printLine(name + " {")
	}
	printing.increaseIndent()
	printing.noNewLineAfterNode--
}

func (printing *Printing) printNodeEnd() {
	printing.decreaseIndent()
	printing.printIndent()
	if printing.colored {
		printing.print(color.RedString("}"))
	} else {
		printing.print("}")
	}
	if printing.noNewLineAfterNode == 0 {
		printing.printNewLine()
	}
}

func (printing *Printing) printFormatted(message string, arguments ...interface{}) {
	printing.buffer.WriteString(fmt.Sprintf(message, arguments...))
}

func (printing *Printing) printIndent() {
	for count := 0; count < printing.indent; count++ {
		printing.print(prettyPrintIndent)
	}
}

func (printing *Printing) printNewLine() {
	printing.print("\n")
}

func (printing *Printing) printNode(node Node) {
	node.Accept(printing.visitor)
}

func (printing *Printing) printListField(node Node) {
	printing.printIndent()
	printing.printNode(node)
	printing.printNewLine()
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
	if printing.colored {
		printing.print(color.RedString("!"))
	} else {
		printing.print("!")
	}
	printing.printNode(statement.Expression)
}

func (printing *Printing) printTranslationUnit(unit *TranslationUnit) {
	printing.printNodeBegin("TranslationUnit")
	printing.printIndentedStringField("name", unit.Name)
	printing.printIndentedListFieldBegin("imports")
	for _, node := range unit.Imports {
		printing.printListField(node)
	}
	printing.printListFieldEnd()
	printing.printIndentedNodeField("class", unit.Class)
	printing.printNodeEnd()
}

func (printing *Printing) printClassDeclaration(class *ClassDeclaration) {
	printing.printNodeBegin("ClassDeclaration")
	printing.printIndentedStringField("name", class.Name)
	printing.printIndentedListFieldBegin("children")
	for _, node := range class.Children {
		printing.printListField(node)
	}
	printing.printListFieldEnd()
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
	printing.printIndentedNodeField("enumeration", statement.Sequence)
	printing.printIndentedNodeField("body", statement.Body)
	printing.printNodeEnd()
}

func (printing *Printing) printMethodDeclaration(method *MethodDeclaration) {
	printing.printNodeBegin("MethodDeclaration")
	printing.printIndentedNodeField("name", method.Name)
	printing.printIndentedNodeField("returnType", method.Type)
	printing.printIndentedListFieldBegin("parameters")
	for _, parameter := range method.Parameters {
		printing.printListField(parameter)
	}
	printing.printListFieldEnd()
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

func (printing *Printing) printListTypeName(name *ListTypeName) {
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

func (printing *Printing) printSelectExpression(expression *SelectExpression) {
	printing.printNodeBegin("Select")
	printing.printIndentedNodeField("target", expression.Target)
	printing.printIndentedNodeField("selection", expression.Selection)
	printing.printNodeEnd()
}

func (printing *Printing) printListSelectExpression(expression *ListSelectExpression) {
	printing.printNodeBegin("ListSelect")
	printing.printIndentedNodeField("target", expression.Target)
	printing.printIndentedNodeField("index", expression.Index)
	printing.printNodeEnd()
}

func (printing *Printing) printMethodCall(call *MethodCall) {
	printing.printNodeBegin("MethodCall")
	printing.printIndentedNodeField("method", call.Method)
	printing.printIndentedListFieldBegin("arguments")
	for _, argument := range call.Arguments {
		printing.printListField(argument)
	}
	printing.printListFieldEnd()
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
		printing.printListField(child)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printImportStatement(statement *ImportStatement) {
	printing.printNodeBegin("ImportStatement")
	target := statement.Target
	targetString := fmt.Sprintf("{path: %s, moduleName: %s}", target.FilePath(), target.toModuleName())
	printing.printIndentedStringField("target", targetString)
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
