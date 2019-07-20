package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/compiler"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/format"
	"os"
)

var standardFormat = format.Format{
	TabWidth:        2,
		ImproveBranches: true,
		IndentWriter:    format.TabIndentWriter{},
		LineLengthLimit: 100,
		EndOfLine:       format.UnixEndOfLine,
	}

var formatCommand = &cobra.Command{
	Use: "format [-f format] [-o override] [-c] [file]",
	Short: "Formats a source file",
	Long: `Format rewrites a strict source file according to the
standard strict-formatting guidelines.`,
	Run: RunFormat,
}

var (
	formatTargetFile	string
	overrideSourceFile bool
)

func init() {
	formatCommand.Flags().
		StringVarP(&formatTargetFile, "target", "t", "", "path to the output file")

	expectNoError(buildCommand.MarkFlagFilename("target", "strict"))
	formatCommand.Flags().BoolVarP(&overrideSourceFile, "override", "o", true, "compile the generated cpp code")
}

func RunFormat(command *cobra.Command, arguments []string) {
	sourceFile, ok := findSourceFileInArguments(command, arguments)
	if !ok {
		return
	}
	defer sourceFile.Close()
	formatFile(command, sourceFile)
}

func formatFile(command *cobra.Command, sourceFile *os.File) {
	parseResult := compiler.ParseFile(sourceFile)
	parseResult.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{command})
	if parseResult.Error != nil {
		return
	}
	targetFile, err := formattingTargetFile(sourceFile)
	if err != nil {
		command.PrintErrf("Failed to create target file: %s\n", err.Error())
	}
	defer targetFile.Close()

	if err := formatUnitToFile(parseResult.Unit, targetFile); err != nil {
		command.PrintErrf("Failed to write formatted sources: %s\n", err.Error())
	} else {
		command.Println("Successfully formatted file")
	}
}

func formattingTargetFile(sourceFile *os.File) (*os.File, error) {
	if overrideSourceFile {
		return sourceFile, nil
	}
	if formatTargetFile == "" {
		return createNewFile(fmt.Sprintf("formatted.%s", sourceFile.Name()))
	}
	return createNewFile(formatTargetFile)
}

func formatUnitToFile(unit *ast.TranslationUnit, file *os.File) (err error) {
	buffer := format.NewStringWriter()
	formatUnit(unit, buffer)
	_, err =  file.WriteString(buffer.String())
	return
}

func formatUnit(unit *ast.TranslationUnit, writer format.Writer) {
	factory := &format.PrettyPrinterFactory{
		Format: standardFormat,
		Unit: unit,
		Writer: writer,
	}
	factory.NewPrettyPrinter().Print()
}
