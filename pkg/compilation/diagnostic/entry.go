package diagnostic

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
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
	Column    source2.Offset
	LineIndex source2.LineIndex
}

func (entry Entry) PrintColored(printer Printer) {
	line := entry.Position.LineIndex

	highlight := entry.Kind.Color.SprintFunc()
	underscore := underlinedColor.SprintFunc()

	PrintFormatted("%s at line %s in %s:  ",
		highlight(entry.Kind.Name), highlight(line), highlight(entry.UnitName))

	Print(underscore(entry.Source))
	PrintFormatted("\n  => %s\n", highlight(entry.Message))
}
