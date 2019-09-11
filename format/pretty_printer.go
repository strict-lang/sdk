package format

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"unicode/utf8"
)

type PrettyPrinter struct {
	format     Format
	writer     Writer
	indent     Indent
	unit       *ast.TranslationUnit
	lineLength int
	astVisitor *ast.Visitor
	sawReturn  bool
}

type PrettyPrinterFactory struct {
	Format Format
	Writer Writer
	Unit   *ast.TranslationUnit
}

func (factory *PrettyPrinterFactory) NewPrettyPrinter() *PrettyPrinter {
	printer := &PrettyPrinter{
		format:     factory.Format,
		writer:     factory.Writer,
		unit:       factory.Unit,
		indent:     Indent{},
		astVisitor: ast.NewEmptyVisitor(),
	}
	printer.registerAstVisitors()
	return printer
}

func (printer *PrettyPrinter) Print() {
	printer.printNode(printer.unit)
}

func (printer *PrettyPrinter) append(text string) {
	printer.writer.Write(text)
	printer.lineLength += utf8.RuneCount([]byte(text))
}

func (printer *PrettyPrinter) appendFormatted(text string, arguments ...interface{}) {
	printer.append(fmt.Sprintf(text, arguments...))
}

func (printer *PrettyPrinter) appendRune(value rune) {
	printer.writer.WriteRune(value)
	printer.lineLength++
}

func (printer *PrettyPrinter) appendIndent() {
	printer.format.IndentWriter.Write(
		printer.indent, printer.writer)
}

func (printer *PrettyPrinter) appendLineBreak() {
	printer.lineLength = 0
	printer.append(printer.format.EndOfLine)
}

func (printer *PrettyPrinter) swapWriter(writer Writer) Writer {
	oldWriter := printer.writer
	printer.writer = writer
	return oldWriter
}

func (printer *PrettyPrinter) setWriter(writer Writer) {
	printer.writer = writer
}

func (printer *PrettyPrinter) printNode(node ast.Node) {
	node.Accept(printer.astVisitor)
}

func (printer *PrettyPrinter) printTranslationUnit(unit *ast.TranslationUnit) {
	// FIXME
}

func (printer *PrettyPrinter) registerAstVisitors() {
	printer.astVisitor.VisitMethodDeclaration = printer.printMethod
	printer.astVisitor.VisitMethodCall = printer.printMethodCall
	printer.astVisitor.VisitIdentifier = printer.printIdentifier
	printer.astVisitor.VisitNumberLiteral = printer.printNumberLiteral
	printer.astVisitor.VisitStringLiteral = printer.printStringLiteral
	printer.astVisitor.VisitBlockStatement = printer.printBlockStatement
	printer.astVisitor.VisitYieldStatement = printer.printYieldStatement
	printer.astVisitor.VisitTranslationUnit = printer.printTranslationUnit
	printer.astVisitor.VisitUnaryExpression = printer.printUnaryExpression
	printer.astVisitor.VisitReturnStatement = printer.printReturnStatement
	printer.astVisitor.VisitAssignStatement = printer.printAssignStatement
	printer.astVisitor.VisitBinaryExpression = printer.printBinaryExpression
	printer.astVisitor.VisitSelectorExpression = printer.printSelectorExpression
	printer.astVisitor.VisitRangedLoopStatement = printer.printRangedLoopStatement
	printer.astVisitor.VisitExpressionStatement = printer.printExpressionStatement
	printer.astVisitor.VisitConditionalStatement = printer.printConditionalStatement
	printer.astVisitor.VisitForEachLoopStatement = printer.printForEachLoopStatement
}
