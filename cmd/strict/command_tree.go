package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/pkg/compiler"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"os"
)

var treeCommand = &cobra.Command{
	Use:   "tree [-c] [file]",
	Short: "Prints the files AST",
	Long:  `Tree parses a file and prints its Abstract Syntax Tree`,
	Run:   runTreeCommand,
}

func runTreeCommand(command *cobra.Command, arguments []string) {
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
	pretty.PrintColored(parseResult.Unit)
}
