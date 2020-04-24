package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
)

type Result struct {
	Error           error
	Diagnostics     *diagnostic.Diagnostics
	TranslationUnit *tree.TranslationUnit
	LineMap *linemap.LineMap
}

func Parse(name string, reader input.Reader) Result {
	diagnosticBag := diagnostic.NewBag()
	tokenReader := lexical.NewScanning(reader)
	parserFactory := NewDefaultFactory().
		WithTokenStream(tokenReader).
		WithDiagnosticBag(diagnosticBag).
		WithUnitName(name)

	unit, err := parserFactory.NewParser().Parse()
	lineMap := tokenReader.NewLineMap()
	offsetConverter := lineMap.PositionAtOffset
	diagnostics := diagnosticBag.CreateDiagnostics(offsetConverter)
	return Result{
		Error:           err,
		TranslationUnit: unit,
		Diagnostics:     diagnostics,
		LineMap: lineMap,
	}
}

func ParseString(name string, text string) Result {
	return Parse(name, input.NewStringReader(text))
}
