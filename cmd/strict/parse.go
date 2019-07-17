package main

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/parser"
	"gitlab.com/strict-lang/sdk/compiler/scanner"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/compiler/source/linemap"
)

func logDiagnostics(recorder *diagnostic.Recorder, converter diagnostic.OffsetToPositionConverter) {
	diagnostics := recorder.CreateDiagnostics(converter)
	diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
}

type parseResult struct {
	parsedUnit *ast.TranslationUnit
	lines *linemap.Linemap
	err error
}
func parseFile(unitName string, recorder *diagnostic.Recorder, reader source.Reader) parseResult {
	tokenSource := scanner.NewDiagnosticScanner(reader, recorder)
	factory := &parser.Factory{
		TokenReader: tokenSource,
		UnitName:    unitName,
		Recorder:    recorder,
	}
	unit, err := factory.NewParser().ParseTranslationUnit()
	return parseResult{
		parsedUnit: unit,
		lines: tokenSource.CreateLinemap(),
		err: err,
	}
}
