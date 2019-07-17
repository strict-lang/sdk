package main

import (
	"bufio"
	"github.com/urfave/cli"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/format"
	"log"
	"os"
)

func formatCode(context *cli.Context) error {
	if context.NArg() < 1 {
		return cli.NewExitError(context.Command.ArgsUsage, StatusInvalidArguments)
	}
	filename := context.Args()[0]
	targetFilename := findFormattingTargetFile(filename, context)
	if err := compileToDirectory(filename, targetFilename); err != nil {
		return err
	}
	log.Printf("formatted %s\n", filename)
	return nil
}

func findFormattingTargetFile(filename string, context *cli.Context) string {
	if context.Bool("override") {
		return filename
	}
	target := context.String("target")
	if target == "" {
		return filename
	}
	return target
}

func formatToTargetFile(filename string, targetFilename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return ErrNoSuchFile
	}
	defer file.Close()
	unitName, err := ParseUnitName(filename)
	if err != nil {
		return err
	}
	return formatFileWithUnitToTargetFile(unitName, file, targetFilename)
}

func formatFileWithUnitToTargetFile(unitName string, file *os.File, targetFilename string) error {
	recorder := diagnostic.NewRecorder()
	reader := source.NewStreamReader(bufio.NewReader(file))
	result := parseFile(unitName, recorder, reader)
	defer logDiagnostics(recorder, result.lines.PositionAtOffset)
	if result.err != nil {
		return reportFailedCompilation(result.err)
	}
	return writeFormattedUnitToFile(result.parsedUnit, targetFilename)
}

func writeFormattedUnitToFile(unit *ast.TranslationUnit, filename string) (err error) {
	buffer := format.NewStringWriter()
	createPrettyPrinter(unit, buffer).Print()
	file, err := createNewFile(filename, "")
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write([]byte(buffer.String()))
	return
}

func createPrettyPrinter(unit *ast.TranslationUnit, writer format.Writer) *format.PrettyPrinter {
	factory := &format.PrettyPrinterFactory{
		Format: format.Format{
			EndOfLine: format.UnixEndOfLine,
			LineLengthLimit: 80,
			IndentWriter: format.ComplexSpaceIndentWriter{
				SpacesPerLevel: 2,
				ContinuousIndent: 4,
			},
			ImproveBranches: true,
			TabWidth: 2,
		},
		Writer: writer,
		Unit: unit,
	}
	return factory.NewPrettyPrinter()
}
