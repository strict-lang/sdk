package compilation

import (
	"gitlab.com/strict-lang/sdk/compilation/backend"
	"gitlab.com/strict-lang/sdk/compilation/backend/headerfile"
	"gitlab.com/strict-lang/sdk/compilation/backend/sourcefile"
	"gitlab.com/strict-lang/sdk/compilation/diagnostic"
	"gitlab.com/strict-lang/sdk/compilation/parsing"
	"gitlab.com/strict-lang/sdk/compilation/scanning"
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
)

type Compilation struct {
	Source        Source
	Name          string
	TargetArduino bool
}

type Result struct {
	UnitName       string
	GeneratedFiles []Generated
	Diagnostics    *diagnostic.Diagnostics
	Error          error
}

type Generated struct {
	FileName string
	Bytes    []byte
}

// ParseResult contains the result of a Parsing.
type ParseResult struct {
	Unit        *syntaxtree.TranslationUnit
	Diagnostics *diagnostic.Diagnostics
	Error       error
}

func (compilation *Compilation) Compile() Result {
	parseResult := compilation.parse()
	if parseResult.Error != nil {
		return Result{
			GeneratedFiles: []Generated{},
			Diagnostics:    parseResult.Diagnostics,
			Error:          parseResult.Error,
			UnitName:       compilation.Name,
		}
	}
	return Result{
		GeneratedFiles: compilation.generateCppFile(parseResult.Unit),
		Diagnostics:    parseResult.Diagnostics,
		Error:          nil,
		UnitName:       parseResult.Unit.Name,
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

func (compilation *Compilation) generateCppFile(unit *syntaxtree.TranslationUnit) []Generated {
	generated := make(chan Generated)
	go func() {
		generated <- compilation.generateHeaderFile(unit)
	}()
	go func() {
		generated <- compilation.generateSourceFile(unit)
	}()
	return []Generated{
		<-generated,
		<-generated,
	}
}

func (compilation *Compilation) generateHeaderFile(unit *syntaxtree.TranslationUnit) Generated {
	generation := backend.NewGenerationWithExtension(unit, headerfile.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".h",
		Bytes:    []byte(generation.Generate()),
	}
}

func (compilation *Compilation) generateSourceFile(unit *syntaxtree.TranslationUnit) Generated {
	generation := backend.NewGenerationWithExtension(unit, sourcefile.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".cc",
		Bytes:    []byte(generation.Generate()),
	}
}
