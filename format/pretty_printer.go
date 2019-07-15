package format

import "gitlab.com/strict-lang/sdk/compiler/ast"

type PrettyPrinter struct {
	format     Format
	writer     Writer
	indent     Indent
	unit       *ast.TranslationUnit
	lineLength int
	astVisitor *ast.Visitor
	sawReturn bool
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
}

func (printer *PrettyPrinter) appendRune(value rune) {
	printer.writer.WriteRune(value)
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
func (printer *PrettyPrinter) registerAstVisitors() {

}
