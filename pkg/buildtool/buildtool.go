package buildtool

import (
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/report"
	"log"
)

func TranslateDiagnostics(diagnostics *diagnostic.Diagnostics) []report.Diagnostic {
	var translated []report.Diagnostic
	for _, entry := range diagnostics.ListEntries() {
		translated = append(translated, translateDiagnosticEntry(entry))
	}
	if translated == nil {
		return []report.Diagnostic{}
	}
	return translated
}

func translateDiagnosticEntry(entry diagnostic.Entry) report.Diagnostic {
	return report.Diagnostic{
		Kind:    translateDiagnosticKind(entry.Kind),
		Message: entry.Message,
		TextRange: report.TextRange{
			Text: entry.Source,
			Range: report.PositionRange{
				BeginPosition: translatePosition(entry.Position.Begin),
				EndPosition:   translatePosition(entry.Position.End),
			},
			File: entry.UnitName,
		},
	}
}

func translatePosition(position input.Position) report.Position {
	return report.Position{
		Line:   int(position.Line.Index),
		Column: int(position.Column),
		Offset: int(position.Offset),
	}
}

func translateDiagnosticKind(kind *diagnostic.Kind) report.DiagnosticKind {
	if kind == nil {
		log.Print("diagnostic-kind is nil")
		return report.DiagnosticError
	}
	switch kind.Name {
	case diagnostic.Error.Name:
		return report.DiagnosticError
	case diagnostic.Info.Name:
		return report.DiagnosticInfo
	case diagnostic.Warning.Name:
		return report.DiagnosticWarning
	default:
		return report.DiagnosticError
	}
}
