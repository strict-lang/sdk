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
	Column    input.Offset
	LineIndex input.LineIndex
}

func (entry Entry) PrintColored(printer Printer) {
	line := entry.Position.LineIndex

	highlight := entry.Kind.Color.SprintFunc()
	underscore := underlinedColor.SprintFunc()

	printer.PrintFormatted("%s at line %s in %s",
		highlight(entry.Kind.Name), highlight(line), highlight(entry.UnitName))

	printer.PrintFormatted("\n  %s %s\n", highlight("at"), underscore(entry.Source))
	entry.printError(printer)
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
	printer.output.PrintLine("\nThis error typically occurs when:")
	printer.output.PrintLine("  " + reason)
}

func (printer *errorPrinter) printMultipleCommonReasons(reasons []string) {
	printer.output.PrintLine("This error is typically a result of one or more of the following:")
	for _, reason := range reasons {
		printer.output.PrintFormatted(" - %s\n", reason)
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
