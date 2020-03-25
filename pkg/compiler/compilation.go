package compiler

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/backend"
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"strict.dev/sdk/pkg/compiler/grammar/syntax"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	isolates "strict.dev/sdk/pkg/compiler/isolate"
	"strict.dev/sdk/pkg/compiler/lowering"
	"strict.dev/sdk/pkg/compiler/pass"
)

type Compilation struct {
	Source  Source
	Name    string
	Backend string
}

type Result struct {
	UnitName       string
	GeneratedFiles []backend.GeneratedFile
	Diagnostics    *diagnostic.Diagnostics
	Error          error
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

func (compilation *Compilation) Compile() Result {
	parseResult := compilation.parse()
	if parseResult.Error != nil {
		return Result{
			GeneratedFiles: []backend.GeneratedFile{},
			Diagnostics:    parseResult.Diagnostics,
			Error:          parseResult.Error,
			UnitName:       compilation.Name,
		}
	}
	compilation.Lower(parseResult.TranslationUnit)
	return Result{
		GeneratedFiles: compilation.generateOutput(parseResult.TranslationUnit),
		Diagnostics:    parseResult.Diagnostics,
		Error:          nil,
		UnitName:       parseResult.TranslationUnit.Name,
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
	unit *tree.TranslationUnit) []backend.GeneratedFile {

	output, _ := compilation.invokeBackend(backend.Input{
		Unit:        unit,
		Diagnostics: diagnostic.NewBag(),
	})
	return output.GeneratedFiles
}

func (compilation *Compilation) invokeBackend(
	input backend.Input) (backend.Output, error) {

	isolate := isolates.SingleThreaded()
	backendId := compilation.Backend
	chosenBackend, ok := backend.LookupInIsolate(isolate, backendId)
	if ok {
		return chosenBackend.Generate(input)
	}
	return backend.Output{}, fmt.Errorf("no such backend: %s", backendId)
}
