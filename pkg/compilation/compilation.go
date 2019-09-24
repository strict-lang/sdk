package compilation

import (
	backend2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	arduino2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend/arduino"
	headerfile2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend/headerfile"
	sourcefile2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend/sourcefile"
	diagnostic2 "gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	parsing2 "gitlab.com/strict-lang/sdk/pkg/compilation/parsing"
	scanning2 "gitlab.com/strict-lang/sdk/pkg/compilation/scanning"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Compilation struct {
	Source        Source
	Name          string
	TargetArduino bool
}

type Result struct {
	UnitName       string
	GeneratedFiles []Generated
	Diagnostics    *diagnostic2.Diagnostics
	Error          error
}

type Generated struct {
	FileName string
	Bytes    []byte
}

// ParseResult contains the result of a Parsing.
type ParseResult struct {
	Unit        *syntaxtree2.TranslationUnit
	Diagnostics *diagnostic2.Diagnostics
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
	diagnosticBag := diagnostic2.NewBag()
	sourceReader := newSourceReader()
	tokenReader := scanning2.NewScanning(sourceReader)
	parserFactory := parsing2.NewDefaultFactory().
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

func (compilation *Compilation) generateOutput(unit *syntaxtree2.TranslationUnit) []Generated {
	if compilation.TargetArduino {
		return []Generated{compilation.generateArduinoFile(unit)}
	}
	return compilation.generateCppFile(unit)
}

func (compilation *Compilation) generateCppFile(unit *syntaxtree2.TranslationUnit) []Generated {
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

func (compilation *Compilation) generateArduinoFile(unit *syntaxtree2.TranslationUnit) Generated {
	generation := backend2.NewGenerationWithExtension(unit, arduino2.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".ino",
		Bytes:    []byte(generation.Generate()),
	}
}

func (compilation *Compilation) generateHeaderFile(unit *syntaxtree2.TranslationUnit) Generated {
	generation := backend2.NewGenerationWithExtension(unit, headerfile2.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".h",
		Bytes:    []byte(generation.Generate()),
	}
}

func (compilation *Compilation) generateSourceFile(unit *syntaxtree2.TranslationUnit) Generated {
	generation := backend2.NewGenerationWithExtension(unit, sourcefile2.NewGeneration())
	return Generated{
		FileName: compilation.Name + ".cc",
		Bytes:    []byte(generation.Generate()),
	}
}
