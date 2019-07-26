package compiler

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/codegen"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/parser"
	"gitlab.com/strict-lang/sdk/compiler/scanner"
	"os"
)

type Compilation struct {
	Source Source
	Name string
}

type CompilationResult struct {
	UnitName    string
	Generated   []byte
	Diagnostics *diagnostic.Diagnostics
	Error       error
}

type ParseResult struct {
	Unit        *ast.TranslationUnit
	Diagnostics *diagnostic.Diagnostics
	Error       error
}

func (compilation *Compilation) Parse() ParseResult {
	recorder := diagnostic.NewRecorder()
	sourceReader := compilation.Source.newSourceReader()
	tokenReader := scanner.NewScanner(sourceReader)
	parserFactory := parser.NewDefaultFactory().
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

func (compilation *Compilation) Run() CompilationResult {
	parseResult := compilation.Parse()
	if parseResult.Error != nil {
		return CompilationResult{
			Generated:   []byte{},
			Diagnostics: parseResult.Diagnostics,
			Error:       parseResult.Error,
			UnitName:    compilation.Name,
		}
	}
	generated := codegen.NewCodeGenerator(parseResult.Unit).Generate()
	return CompilationResult{
		Generated:   []byte(generated),
		Diagnostics: parseResult.Diagnostics,
		Error:       nil,
		UnitName:    parseResult.Unit.Name(),
	}
}

func CompileFile(name string, file *os.File) CompilationResult {
	compilation := &Compilation{
		Source: &FileSource{File: file},
		Name: name,
	}
	return compilation.Run()
}

func ParseFile(name string, file *os.File) ParseResult {
	compilation := &Compilation{
		Source: &FileSource{File: file},
		Name: name,
	}
	return compilation.Parse()
}

func CompileString(name string, value string) CompilationResult {
	compilation := &Compilation{
		Source: &InMemorySource{Source: value},
		Name: name,
	}
	return compilation.Run()
}
