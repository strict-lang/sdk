package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/pkg/compilation"
	"gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	"os"
)

var astCommand = &cobra.Command{
	Use:   "syntaxtree [-c] [file]",
	Short: "Prints the files AST",
	Long:  `Ast parses a file and prints its Abstract Syntax Tree`,
	Run:   RunAst,
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
	parseResult := compilation.ParseFile("formatted", sourceFile)
	parseResult.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{command})
	if parseResult.Error != nil {
		return
	}
	syntaxtree.PrintColored(parseResult.Unit)
}
