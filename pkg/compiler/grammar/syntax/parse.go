package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type Result struct {
	Error           error
	Diagnostics     *diagnostic.Diagnostics
	TranslationUnit *tree.TranslationUnit
}

func Parse(name string, reader input.Reader) Result {
	diagnosticBag := diagnostic.NewBag()
	tokenReader := lexical.NewScanning(reader)
	parserFactory := NewDefaultFactory().
		WithTokenStream(tokenReader).
		WithDiagnosticBag(diagnosticBag).
		WithUnitName(name)

	unit, err := parserFactory.NewParser().Parse()
	offsetConverter := tokenReader.NewLineMap().PositionAtOffset
	diagnostics := diagnosticBag.CreateDiagnostics(offsetConverter)
	return Result{
		Error:           err,
		TranslationUnit: unit,
		Diagnostics:     diagnostics,
	}
}

func ParseString(name string, text string) Result {
	return Parse(name, input.NewStringReader(text))
}
