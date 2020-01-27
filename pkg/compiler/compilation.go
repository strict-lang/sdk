package compiler

import (
	"strict.dev/sdk/pkg/compiler/backend"
	"strict.dev/sdk/pkg/compiler/backend/arduino"
	"strict.dev/sdk/pkg/compiler/backend/headerfile"
	"strict.dev/sdk/pkg/compiler/backend/sourcefile"
	"strict.dev/sdk/pkg/compiler/backend/testfile"
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
		GeneratedFiles: compilation.generateOutput(parseResult.TranslationUnit),
		Diagnostics:    parseResult.Diagnostics,
		Error:          nil,
		UnitName:       parseResult.TranslationUnit.Name,
	}
}

func (compilation *Compilation) parse() syntax.Result {
	return syntax.Parse(compilation.Name, compilation.Source.newSourceReader())
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
	var output []Generated
	if isContainingTestDefinitions(unit) {
		output = append(output, compilation.generateTestFile(unit))
	}
	output = append(output, <-generated)
	output = append(output, <-generated)
	return output
}

func isContainingTestDefinitions(node tree.Node) bool {
	counter := tree.NewCounter()
	counter.Count(node)
	return counter.ValueFor(tree.TestStatementNodeKind) != 0
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

func (compilation *Compilation) generateTestFile(unit *tree.TranslationUnit) Generated {
	generation := backend.NewGenerationWithExtension(unit, testfile.NewGeneration())
	return Generated{
		FileName: compilation.Name + "_test.cc",
		Bytes:    []byte(generation.Generate()),
	}
}
