package diagnostic

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type Entry struct {
	Kind     *Kind
	Stage    *Stage
	Source   string
	Message  string
	UnitName string
	Position Position
	Error    *RichError
}

type Position struct {
	Begin input.Position
	End   input.Position
}

func (position Position) isSpanningMultipleLines() bool {
	return position.Begin.Line.Index != position.End.Line.Index
}

func (position Position) endColumn() input.Offset {
	if position.isSpanningMultipleLines() {
		return position.Begin.Line.Length
	}
	return position.End.Column
}

func (position Position) selectedSource() (begin, end input.Offset) {
	begin = position.Begin.Column
	end = position.endColumn() - 1
	return
}

func (entry Entry) PrintColored(printer Printer) {
	line := entry.Position.Begin.Line.Index

	highlight := entry.Kind.Color.SprintFunc()
	underscore := underlinedColor.SprintFunc()

	printer.PrintFormatted("%s at line %s in %s %s %s\n",
		highlight(entry.Kind.Name),
		highlight(line),
		highlight(entry.UnitName),
		highlight("at"),
		underscore(entry.Source))

	entry.printError(printer)
}

func (entry Entry) formatSource() string {
	begin, end := entry.Position.selectedSource()
	full := entry.Position.Begin.Line.Text
	text := full[begin:end]
	prefix := full[0:begin]
	suffix := full[end:]
	return prefix + "<" + text + ">" + suffix
}

func (entry Entry) printError(printer Printer) {
	errorPrinter := errorPrinter{
		output: printer,
		color:  entry.Kind.Color.SprintFunc(),
	}
	errorPrinter.printRichError(entry.Error)
}

type errorPrinter struct {
	output Printer
	color  func(...interface{}) string
}

func (printer *errorPrinter) printRichError(error *RichError) {
	printer.output.Print("\t")
	error.Error.Accept(printer)
	printer.printCommonReasons(error)
}

func (printer *errorPrinter) printCommonReasons(error *RichError) {
	if len(error.CommonReasons) == 0 {
		return
	}
	if len(error.CommonReasons) == 1 {
		printer.printSingleCommonReason(error.CommonReasons[0])
	} else {
		printer.printMultipleCommonReasons(error.CommonReasons)
	}
}

func (printer *errorPrinter) printSingleCommonReason(reason string) {
	printer.output.PrintLine("\n\tThis error typically occurs when:")
	printer.output.PrintLine("\t  " + reason)
}

func (printer *errorPrinter) printMultipleCommonReasons(reasons []string) {
	printer.output.PrintLine("\n\tThis error is typically a result of one or more of the following:")
	for _, reason := range reasons {
		printer.output.PrintFormatted("\t - %s\n", reason)
	}
}

func (printer *errorPrinter) VisitUnexpectedToken(error *UnexpectedTokenError) {
	printer.output.PrintFormatted(
		"Expected %s but got %s",
		printer.color(error.Expected),
		printer.color(error.Received))
}

func (printer *errorPrinter) VisitInvalidStatement(error *InvalidStatementError) {
	printer.output.PrintFormatted("Could not parse %s", printer.color(error.Kind.Name()))
}

func (printer *errorPrinter) VisitInvalidIndentation(error *InvalidIndentationError) {
	printer.output.PrintFormatted("Expected indentation of %s but got %s",
		printer.color(error.Expected),
		printer.color(error.Received))
}

func (printer *errorPrinter) VisitSpecificError(error *SpecificError) {
	printer.output.Print(error.Message)
}
