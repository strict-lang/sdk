package compiler

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend/arduino"
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend/headerfile"
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend/sourcefile"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/syntax"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
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
	Unit        *tree.TranslationUnit
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
		GeneratedFiles: compilation.generateOutput(parseResult.Unit),
		Diagnostics:    parseResult.Diagnostics,
		Error:          nil,
		UnitName:       parseResult.Unit.Name,
	}
}

func (compilation *Compilation) parse() ParseResult {
	diagnosticBag := diagnostic.NewBag()
	sourceReader := compilation.Source.newSourceReader()
	tokenReader := lexical.NewScanning(sourceReader)
	parserFactory := syntax.NewDefaultFactory().
		WithTokenStream(tokenReader).
		WithDiagnosticBag(diagnosticBag).
		WithUnitName(compilation.Name)

	unit, err := parserFactory.NewParser().Parse()
	offsetConverter := tokenReader.NewLineMap().PositionAtOffset
	diagnostics := diagnosticBag.CreateDiagnostics(offsetConverter)
	return ParseResult{
		Unit:        unit,
		Diagnostics: diagnostics,
		Error:       err,
	}
}

func (compilation *Compilation) generateOutput(unit *tree.TranslationUnit) []Generated {
	if compilation.TargetArduino {
		return []Generated{compilation.generateArduinoFile(unit)}
	}
	return compilation.generateCppFile(unit)
}

func (compilation *Compilation) generateCppFile(unit *tree.TranslationUnit) []Generated {
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

func (compilation *Compilation) generateArduinoFile(unit *tree.TranslationUnit) Generated {
	generation := backend.NewGenerationWithExtension(unit, arduino.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".ino",
		Bytes:    []byte(generation.Generate()),
	}
}

func (compilation *Compilation) generateHeaderFile(unit *tree.TranslationUnit) Generated {
	generation := backend.NewGenerationWithExtension(unit, headerfile.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".h",
		Bytes:    []byte(generation.Generate()),
	}
}

func (compilation *Compilation) generateSourceFile(unit *tree.TranslationUnit) Generated {
	generation := backend.NewGenerationWithExtension(unit, sourcefile.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".cc",
		Bytes:    []byte(generation.Generate()),
	}
}
