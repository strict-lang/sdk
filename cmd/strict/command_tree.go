package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/pkg/compiler"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree/pretty"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	"gitlab.com/strict-lang/sdk/pkg/compiler/lowering"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
	"os"
	"strings"
)

var treeCommand = &cobra.Command{
	Use:   "tree [-c] [file]",
	Short: "Prints the files AST",
	Long:  `Tree parses a file and prints its Abstract Syntax Tree`,
	Run:   runTreeCommand,
}

func runTreeCommand(command *cobra.Command, arguments []string) {
	if sourceFile, ok := findSourceFileInArguments(command, arguments); ok {
		defer sourceFile.Close()
		parseAndPrintAst(command, sourceFile)
	}
}

func createUnitNameFromFileName(name string) string {
	return strings.TrimSuffix(name, ".strict")
}

func parseAndPrintAst(command *cobra.Command, sourceFile *os.File) {
	parseResult := compiler.ParseFile(sourceFile.Name(), sourceFile)
	parseResult.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{command})
	if parseResult.Error != nil {
		return
	}
	analyseAndLowerUnit(parseResult.Unit)
	pretty.PrintColored(parseResult.Unit)
}

func analyseAndLowerUnit(unit *tree.TranslationUnit) {
	pass, _ := passes.NewExecution(lowering.LetBindingLoweringPassId, &passes.Context{
		Unit:       unit,
		Diagnostic: diagnostic.NewBag(),
		Isolate:    isolate.SingleThreaded(),
	})
	_ = pass.Run()
}
