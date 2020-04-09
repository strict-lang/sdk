package pretty

import (
	"fmt"
	"github.com/fatih/color"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"strings"
)

const prettyPrintIndent = "  "

type Printing struct {
	buffer             strings.Builder
	indent             int
	visitor            tree.Visitor
	colored            bool
	noNewLineAfterNode int
}

func newPrinting() *Printing {
	printing := &Printing{}
	visitor := &tree.DelegatingVisitor{
		ParameterVisitor:              printing.printParameter,
		CallExpressionVisitor:         printing.printCallExpression,
		CallArgumentVisitor:           printing.printCallArgument,
		IdentifierVisitor:             printing.printIdentifier,
		ListTypeNameVisitor:           printing.printListTypeName,
		TestStatementVisitor:          printing.printTestStatement,
		StringLiteralVisitor:          printing.printStringLiteral,
		NumberLiteralVisitor:          printing.printNumberLiteral,
		EmptyStatementVisitor:         printing.printEmptyStatement,
		YieldStatementVisitor:         printing.printYieldStatement,
		BlockStatementVisitor:         printing.printBlockStatement,
		AssertStatementVisitor:        printing.printAssertStatement,
		UnaryExpressionVisitor:        printing.printUnaryExpression,
		ImportStatementVisitor:        printing.printImportStatement,
		AssignStatementVisitor:        printing.printAssignStatement,
		ReturnStatementVisitor:        printing.printReturnStatement,
		TranslationUnitVisitor:        printing.printTranslationUnit,
		CreateExpressionVisitor:       printing.printCreateExpression,
		InvalidStatementVisitor:       printing.printInvalidStatement,
		FieldDeclarationVisitor:       printing.printFieldDeclaration,
		ClassDeclarationVisitor:       printing.printClassDeclaration,
		GenericTypeNameVisitor:        printing.printGenericTypeName,
		ConcreteTypeNameVisitor:       printing.printConcreteTypeName,
		BinaryExpressionVisitor:       printing.printBinaryExpression,
		MethodDeclarationVisitor:      printing.printMethodDeclaration,
		FieldSelectExpressionVisitor:  printing.printFieldSelectExpression,
		PostfixExpressionVisitor:      printing.printPostfixExpression,
		RangedLoopStatementVisitor:    printing.printRangedLoopStatement,
		ExpressionStatementVisitor:    printing.printExpressionStatement,
		ForEachLoopStatementVisitor:   printing.printForEachLoopStatement,
		ConditionalStatementVisitor:   printing.printConditionalStatement,
		LetBindingVisitor:             printing.printLetBinding,
		ListSelectExpressionVisitor:   printing.printListSelectExpression,
		ConstructorDeclarationVisitor: printing.printConstructorDeclaration,
		WildcardNodeVisitor:           printing.printWildcardNode,
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
func Print(node tree.Node) {
	printing := newPrinting()
	printing.printNode(node)
	fmt.Println(printing.buffer.String())
}

func Format(node tree.Node) string {
	printing := newPrinting()
	printing.printNode(node)
	return printing.buffer.String()
}

// PrintColored prints a pretty representation of the AST node, that uses
// colors to highlight the output in a way that makes it easier to read.
func PrintColored(node tree.Node) {
	printing := newColoredPrinting()
	printing.printNode(node)
	fmt.Println(printing.buffer.String())
}

func FormatColored(node tree.Node) string {
	printing := newColoredPrinting()
	printing.printNode(node)
	return printing.buffer.String()
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

func (printing *Printing) printIndentedNodeField(name string, node tree.Node) {
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

func (printing *Printing) printNode(node tree.Node) {
	if node == nil {
		if printing.colored {
			printing.print(color.RedString("!!!nil"))
		} else {
			printing.print("!!!nil")
		}
		return
	}
	node.Accept(printing.visitor)
}

func (printing *Printing) printResolvedType(expression tree.Expression) {
	if resolvedType, ok := expression.ResolvedType(); ok {
		printing.printIndentedStringField("type", resolvedType.Name())
	}
}

func (printing *Printing) printListField(node tree.Node) {
	printing.printIndent()
	printing.printNode(node)
	printing.printNewLine()
}

func (printing *Printing) printBinaryExpression(expression *tree.BinaryExpression) {
	printing.printNodeBegin("BinaryExpression")
	printing.printIndentedStringField("operator", expression.Operator.String())
	printing.printIndentedNodeField("leftOperand", expression.LeftOperand)
	printing.printIndentedNodeField("rightOperand", expression.RightOperand)
	printing.printResolvedType(expression)
	printing.printNodeEnd()
}

func (printing *Printing) printUnaryExpression(expression *tree.UnaryExpression) {
	printing.printNodeBegin("UnaryExpression")
	printing.printIndentedStringField("operator", expression.Operator.String())
	printing.printIndentedNodeField("operand", expression.Operand)
	printing.printResolvedType(expression)
	printing.printNodeEnd()
}

func (printing *Printing) printExpressionStatement(statement *tree.ExpressionStatement) {
	if printing.colored {
		printing.print(color.RedString("!"))
	} else {
		printing.print("!")
	}
	printing.printNode(statement.Expression)
}

func (printing *Printing) printTranslationUnit(unit *tree.TranslationUnit) {
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

func (printing *Printing) printClassDeclaration(class *tree.ClassDeclaration) {
	printing.printNodeBegin("ClassDeclaration")
	printing.printIndentedStringField("name", class.Name)
	printing.printIndentedListFieldBegin("children")
	for _, node := range class.Children {
		printing.printListField(node)
	}
	printing.printListFieldEnd()
}

func (printing *Printing) printLetBinding(binding *tree.LetBinding) {
	printing.printNodeBegin("LetBinding")
	if len(binding.Names) > 1 {
		printing.printIndentedListFieldBegin("names")
		for _, name := range binding.Names {
			printing.printListField(name)
		}
		printing.printListFieldEnd()
	} else {
		printing.printIndentedNodeField("name", binding.Names[0])
	}
	printing.printResolvedType(binding)
	printing.printIndentedNodeField("value", binding.Expression)
	printing.printNodeEnd()
}

func (printing *Printing) printIdentifier(identifier *tree.Identifier) {
	_, resolved := identifier.ResolvedType()
	if !resolved && !identifier.IsBound() {
		printing.printFormatted("%s", identifier.Value)
		return
	}
	printing.printRichIdentifier(identifier)
}

func (printing *Printing) printRichIdentifier(identifier *tree.Identifier) {
	printing.printNodeBegin("Identifier")
	printing.printIndentedStringField("value", identifier.Value)
	printing.printResolvedType(identifier)
	if identifier.IsBound() {
		printing.printIndentedStringField("boundTo", identifier.Binding().String())
	}
	printing.printNodeEnd()
}

func (printing *Printing) printStringLiteral(literal *tree.StringLiteral) {
	printing.printFormatted("\"%s\"", literal.Value)
}

func (printing *Printing) printNumberLiteral(number *tree.NumberLiteral) {
	printing.print(number.Value)
}

func (printing *Printing) printAssertStatement(statement *tree.AssertStatement) {
	printing.printNodeBegin("AssertStatement")
	printing.printIndentedNodeField("expression", statement.Expression)
	printing.printNodeEnd()
}

func (printing *Printing) printReturnStatement(statement *tree.ReturnStatement) {
	printing.printNodeBegin("ReturnStatement")
	if statement.Value == nil {
		printing.printIndentedStringField("value", "None")
	} else {
		printing.printIndentedNodeField("value", statement.Value)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printYieldStatement(statement *tree.YieldStatement) {
	printing.printNodeBegin("YieldStatement")
	printing.printIndentedNodeField("value", statement.Value)
	printing.printNodeEnd()
}

func (printing *Printing) printConditionalStatement(statement *tree.ConditionalStatement) {
	printing.printNodeBegin("ConditionalStatement")
	printing.printIndentedNodeField("condition", statement.Condition)
	printing.printIndentedNodeField("consequence", statement.Consequence)
	if statement.HasAlternative() {
		printing.printIndentedNodeField("Alternative", statement.Alternative)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printRangedLoopStatement(statement *tree.RangedLoopStatement) {
	printing.printNodeBegin("RangedLoopStatement")
	printing.printIndentedNodeField("valueField", statement.Field)
	printing.printIndentedNodeField("initialValue", statement.Begin)
	printing.printIndentedNodeField("endValue", statement.End)
	printing.printIndentedNodeField("body", statement.Body)
	printing.printNodeEnd()
}

func (printing *Printing) printForEachLoopStatement(statement *tree.ForEachLoopStatement) {
	printing.printNodeBegin("ForEachLoopStatement")
	printing.printIndentedNodeField("field", statement.Field)
	printing.printIndentedNodeField("enumeration", statement.Sequence)
	printing.printIndentedNodeField("body", statement.Body)
	printing.printNodeEnd()
}

func (printing *Printing) printMethodDeclaration(method *tree.MethodDeclaration) {
	printing.printNodeBegin("MethodDeclaration")
	printing.printIndentedNodeField("name", method.Name)
	printing.printIndentedNodeField("returnType", method.Type)
	printing.printParameterList(method.Parameters)
	if method.Body != nil {
		printing.printIndentedNodeField("body", method.Body)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printParameterList(parameters tree.ParameterList) {
	printing.printIndentedListFieldBegin("parameters")
	for _, parameter := range parameters {
		printing.printListField(parameter)
	}
	printing.printListFieldEnd()
}

func (printing *Printing) printWildcardNode(node *tree.WildcardNode) {
	printing.printFieldName("*")
}

func (printing *Printing) printConstructorDeclaration(declaration *tree.ConstructorDeclaration) {
	printing.printNodeBegin("ConstructorDeclaration")
	printing.printParameterList(declaration.Parameters)
	if declaration.Body != nil {
		printing.printIndentedNodeField("body", declaration.Body)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printParameter(parameter *tree.Parameter) {
	printing.printNodeBegin("Parameter")
	printing.printIndentedNodeField("name", parameter.Name)
	printing.printIndentedNodeField("type", parameter.Type)
	printing.printNodeEnd()
}

func (printing *Printing) printGenericTypeName(name *tree.GenericTypeName) {
	printing.print(name.FullName())
}

func (printing *Printing) printListTypeName(name *tree.ListTypeName) {
	printing.print(name.FullName())
}

func (printing *Printing) printConcreteTypeName(name *tree.ConcreteTypeName) {
	printing.print(name.FullName())
}

func (printing *Printing) printInvalidStatement(statement *tree.InvalidStatement) {
	printing.print("!!!INVALID")
}

func (printing *Printing) printEmptyStatement(statement *tree.EmptyStatement) {
	printing.print("-")
}

func (printing *Printing) printPostfixExpression(statement *tree.PostfixExpression) {
	printing.printNodeBegin("PostfixExpression")
	printing.printIndentedStringField("operator", statement.Operator.String())
	printing.printIndentedNodeField("operand", statement.Operand)
	printing.printResolvedType(statement)
	printing.printNodeEnd()
}

func (printing *Printing) printFieldSelectExpression(expression *tree.FieldSelectExpression) {
	printing.printNodeBegin("Select")
	printing.printIndentedNodeField("target", expression.Target)
	printing.printIndentedNodeField("selection", expression.Selection)
	printing.printResolvedType(expression)
	printing.printNodeEnd()
}

func (printing *Printing) printListSelectExpression(expression *tree.ListSelectExpression) {
	printing.printNodeBegin("ListSelect")
	printing.printIndentedNodeField("target", expression.Target)
	printing.printIndentedNodeField("index", expression.Index)
	printing.printResolvedType(expression)
	printing.printNodeEnd()
}

func (printing *Printing) printCallExpression(call *tree.CallExpression) {
	printing.printNodeBegin("CallExpression")
	printing.printResolvedType(call)
	printing.printIndentedNodeField("method", call.Target)
	printing.printIndentedListFieldBegin("arguments")
	for _, argument := range call.Arguments {
		printing.printListField(argument)
	}
	printing.printListFieldEnd()
	printing.printNodeEnd()
}

func (printing *Printing) printCallArgument(argument *tree.CallArgument) {
	printing.printNodeBegin("CallArgument")
	if argument.IsLabeled() {
		printing.printIndentedStringField("label", argument.Label)
	}
	printing.printIndentedNodeField("value", argument.Value)
	printing.printNodeEnd()
}

func (printing *Printing) printFieldDeclaration(field *tree.FieldDeclaration) {
	printing.printNodeBegin("FieldDeclaration")
	printing.printIndentedNodeField("name", field.Name)
	printing.printIndentedNodeField("type", field.TypeName)
	printing.printNodeEnd()
}

func (printing *Printing) printAssignStatement(statement *tree.AssignStatement) {
	printing.printNodeBegin("AssignStatement")
	printing.printIndentedNodeField("target", statement.Target)
	printing.printIndentedNodeField("value", statement.Value)
	printing.printNodeEnd()
}

func (printing *Printing) printTestStatement(statement *tree.TestStatement) {
	printing.printNodeBegin("TestStatement")
	printing.printIndentedStringField("methodName", statement.MethodName)
	printing.printIndentedNodeField("body", statement.Body)
	printing.printNodeEnd()
}

func (printing *Printing) printBlockStatement(statement *tree.StatementBlock) {
	printing.printNodeBegin("StatementBlock")
	for _, child := range statement.Children {
		printing.printListField(child)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printImportStatement(statement *tree.ImportStatement) {
	printing.printNodeBegin("ImportStatement")
	target := statement.Target
	targetString := fmt.Sprintf("{path: %s, moduleName: %s}", target.FilePath(), target.ToModuleName())
	printing.printIndentedStringField("target", targetString)
	if statement.Alias != nil {
		printing.printIndentedNodeField("alias", statement.Alias)
	}
	printing.printNodeEnd()
}

func (printing *Printing) printCreateExpression(expression *tree.CreateExpression) {
	printing.printNodeBegin("CreateExpression")
	printing.printIndentedNodeField("constructor", expression.Call)
	printing.printNodeEnd()
}
