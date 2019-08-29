package compilation

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/parsing"
	"gitlab.com/strict-lang/sdk/compilation/scanning"
	"os"
)

type Compilation struct {
	Source        Source
	Name          string
	TargetArduino bool
}

type Result struct {
	UnitName          string
	Generated         []byte
	GeneratedFileName string
	Diagnostics       *diagnostic.Diagnostics
	Error             error
}

type ParseResult struct {
	Unit        *ast.TranslationUnit
	Diagnostics *diagnostic.Diagnostics
	Error       error
}

func (compilation *Compilation) Parse() ParseResult {
	recorder := diagnostic.NewRecorder()
	sourceReader := compilation.Source.newSourceReader()
	tokenReader := scanning.NewScanning(sourceReader)
	parserFactory := parsing.NewDefaultFactory().
		WithTokenReader(tokenReader).
		WithRecorder(recorder).
		WithUnitName(compilation.Name)

	unit, err := parserFactory.NewParser().ParseTranslationUnit()
	offsetConverter := tokenReader.CreateLinemap().PositionAtOffset
	diagnostics := recorder.CreateDiagnostics(offsetConverter)
	return ParseResult{
		Unit:        unit,
		Diagnostics: diagnostics,
		Error:       err,
	}
}

func (compilation *Compilation) Run() Result {
	parseResult := compilation.Parse()
	if parseResult.Error != nil {
		return Result{
			Generated:   []byte{},
			Diagnostics: parseResult.Diagnostics,
			Error:       parseResult.Error,
			UnitName:    compilation.Name,
		}
	}
	generator := backend.NewGeneration(parseResult.Unit)
	return Result{
		Generated:         []byte(generator.Generate()),
		GeneratedFileName: generator.Filename(),
		Diagnostics:       parseResult.Diagnostics,
		Error:             nil,
		UnitName:          parseResult.Unit.Name(),
	}
}

func CompileFile(name string, file *os.File) Result {
	compilation := &Compilation{
		Source: &FileSource{File: file},
		Name:   name,
	}
	return compilation.Run()
}

func ParseFile(name string, file *os.File) ParseResult {
	compilation := &Compilation{
		Source: &FileSource{File: file},
		Name:   name,
	}
	return compilation.Parse()
}

func CompileString(name string, value string) Result {
	compilation := &Compilation{
		Source: &InMemorySource{Source: value},
		Name:   name,
	}
	return compilation.Run()
}
