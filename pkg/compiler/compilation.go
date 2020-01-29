package compiler

import (
	"strict.dev/sdk/pkg/compiler/backend"
	"strict.dev/sdk/pkg/compiler/backend/cpp"
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"strict.dev/sdk/pkg/compiler/grammar/syntax"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
)

type Compilation struct {
	Source        Source
	Name          string
	TargetArduino bool
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
	return Result{
		GeneratedFiles: compilation.generateOutput(parseResult.TranslationUnit),
		Diagnostics:    parseResult.Diagnostics,
		Error:          nil,
		UnitName:       parseResult.TranslationUnit.Name,
	}
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
	input backend.Input) (backend.Output, error){

	if compilation.TargetArduino {
		return cpp.Generate(input)
	}
	return cpp.Generate(input)
}



