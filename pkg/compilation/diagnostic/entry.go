package diagnostic

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type Entry struct {
	Kind     *Kind
	Stage    *Stage
	Source   string
	Message  string
	UnitName string
	Position Position
}

type Position struct {
	Column    source.Offset
	LineIndex source.LineIndex
}

func (entry Entry) PrintColored(printer Printer) {
	line := entry.Position.LineIndex

	highlight := entry.Kind.Color.SprintFunc()
	underscore := underlinedColor.SprintFunc()

	printer.PrintFormatted("%s at line %s in %s:  ",
		highlight(entry.Kind.Name), highlight(line), highlight(entry.UnitName))

	printer.Print(underscore(entry.Source))
	printer.PrintFormatted("\n  => %s\n", highlight(entry.Message))
}
