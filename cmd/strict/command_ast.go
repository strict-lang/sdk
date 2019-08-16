package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/compiler"
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"os"
)

var astCommand = &cobra.Command{
	Use:   "ast [-c] [file]",
	Short: "Prints the files AST",
	Long: `Ast parses a file and prints its Abstract Syntax Tree`,
	Run: RunAst,
}

func RunAst(command *cobra.Command, arguments []string) {
	sourceFile, ok := findSourceFileInArguments(command, arguments)
	if !ok {
		return
	}
	defer sourceFile.Close()
	parseAndPrintAst(command, sourceFile)
}

func parseAndPrintAst(command *cobra.Command, sourceFile *os.File) {
	parseResult := compiler.ParseFile("formatted", sourceFile)
	parseResult.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{command})
	if parseResult.Error != nil {
		return
	}
	ast.Print(parseResult.Unit)
}
