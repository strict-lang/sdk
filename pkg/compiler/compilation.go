package compiler

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/backend"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/arduino"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/cpp"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/ilasm"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/silk"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
	isolates "github.com/strict-lang/sdk/pkg/compiler/isolate"
	"github.com/strict-lang/sdk/pkg/compiler/lowering"
	"github.com/strict-lang/sdk/pkg/compiler/pass"
	"github.com/strict-lang/sdk/pkg/compiler/report"
	"log"
	"time"
)

type Compilation struct {
	Source  Source
	Name    string
	Backend string

	beginTime   time.Time
	diagnostics *diagnostic.Diagnostics
	err         error
}

type Result struct {
	UnitName       string
	GeneratedFiles []backend.GeneratedFile
	Report         report.Report
	Error          error
	LineMap        *linemap.LineMap
}

type Generated struct {
	FileName string
	Bytes    []byte
}

// ParseResult contains the result of a Parsing.
type ParseResult struct {
	Unit        *tree.TranslationUnit
	Diagnostics *diagnostic.Diagnostics
	Error       error
}

const failure = false
const success = true

func (compilation *Compilation) createReport(success bool) report.Report {
	return report.Report{
		Success: success,
		Time: report.Time{
			Begin:      compilation.beginTime.UnixNano(),
			Completion: compilation.beginTime.UnixNano(),
		},
		Diagnostics: translateDiagnostics(compilation.diagnostics),
	}
}

func translateDiagnostics(diagnostics *diagnostic.Diagnostics) []report.Diagnostic {
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

func (compilation *Compilation) Compile() Result {
	compilation.beginTime = time.Now()
	parseResult := compilation.parse()
	compilation.diagnostics = parseResult.Diagnostics
	if parseResult.Error != nil {
		log.Printf("could not parse input file: %v", parseResult.Error)
		return Result{
			GeneratedFiles: []backend.GeneratedFile{},
			Report:         compilation.createReport(failure),
			Error:          parseResult.Error,
			UnitName:       compilation.Name,
			LineMap:        parseResult.LineMap,
		}
	}
	compilation.Lower(parseResult.TranslationUnit)
	generatedFiles, err := compilation.generateOutput(parseResult.TranslationUnit)

	return Result{
		GeneratedFiles: generatedFiles,
		Report:         compilation.createReport(success),
		Error:          err,
		UnitName:       parseResult.TranslationUnit.Name,
		LineMap:        parseResult.LineMap,
	}
}

func (compilation *Compilation) Lower(unit *tree.TranslationUnit) {
	execution, _ := pass.NewExecution(lowering.LetBindingLoweringPassId, &pass.Context{
		Unit:       unit,
		Diagnostic: diagnostic.NewBag(),
		Isolate:    isolates.SingleThreaded(),
	})
	_ = execution.Run()
}

func (compilation *Compilation) parse() syntax.Result {
	return syntax.Parse(compilation.Name, compilation.Source.newSourceReader())
}

func (compilation *Compilation) generateOutput(
	unit *tree.TranslationUnit) ([]backend.GeneratedFile, error) {

	output, err := compilation.invokeBackend(backend.Input{
		Unit:        unit,
		Diagnostics: diagnostic.NewBag(),
	})
	if err != nil {
		return nil, err
	}
	return output.GeneratedFiles, nil
}

func (compilation *Compilation) invokeBackend(
	input backend.Input) (backend.Output, error) {

	isolate := isolates.SingleThreaded()
	backendId := compilation.Backend
	chosenBackend, ok := backend.LookupInIsolate(isolate, backendId)
	if ok {
		log.Printf("compiling files with %s backend", backendId)
		return chosenBackend.Generate(input)
	}
	log.Printf("could not find backend %s", backendId)
	return backend.Output{}, fmt.Errorf("no such backend: %s", backendId)
}
