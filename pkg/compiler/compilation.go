package compiler

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/buildtool"
	"github.com/strict-lang/sdk/pkg/compiler/backend"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/arduino"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/cpp"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/ilasm"
	_ "github.com/strict-lang/sdk/pkg/compiler/backend/silk"
	"github.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
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
		Diagnostics: buildtool.TranslateDiagnostics(compilation.diagnostics),
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
	execution, _ := pass.NewExecution(lowering.FullLoweringI, &pass.Context{
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
