package compilation

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/parsing"
	"gitlab.com/strict-lang/sdk/compilation/scanning"
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

// ParseResult contains the result of a Parsing.
type ParseResult struct {
	Unit        *ast.TranslationUnit
	Diagnostics *diagnostic.Diagnostics
	Error       error
}

func (compilation *Compilation) Compile() Result {
	parseResult := compilation.parse()
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
		UnitName:          parseResult.Unit.Name,
	}
}

func (compilation *Compilation) parse() ParseResult {
	diagnosticBag := diagnostic.NewBag()
	sourceReader := compilation.Source.newSourceReader()
	tokenReader := scanning.NewScanning(sourceReader)
	parserFactory := parsing.NewDefaultFactory().
		WithTokenStream(tokenReader).
		WithDiagnosticBag(diagnosticBag).
		WithUnitName(compilation.Name)

	unit, err := parserFactory.NewParser().ParseTranslationUnit()
	offsetConverter := tokenReader.NewLineMap().PositionAtOffset
	diagnostics := diagnosticBag.CreateDiagnostics(offsetConverter)
	return ParseResult{
		Unit:        unit,
		Diagnostics: diagnostics,
		Error:       err,
	}
}
