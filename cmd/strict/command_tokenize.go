package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/lexical"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

var tokenizeCommand = &cobra.Command{
	Use:   "tokenize [-c] [file]",
	Short: "Scans the file and prints the tokens",
	Long:  `Splits the files characters into tokens and prints them`,
	Run:   RunTokenize,
}

func RunTokenize(command *cobra.Command, arguments []string) {
	sourceFile, ok := findSourceFileInArguments(command, arguments)
	if !ok {
		return
	}
	defer sourceFile.Close()
	scanAndPrintTokens(command, sourceFile)
}

func printNewLineIndent() {
	fmt.Print("  ")
}

func scanAndPrintTokens(command *cobra.Command, sourceFile *os.File) {
	scan := lexical.NewScanning(input.NewStreamReader(sourceFile))
	fmt.Println("Scanned tokens:")
	printNewLineIndent()
	for {
		next := scan.Pull()
		if token.IsEndOfFileToken(next) {
			break
		}
		fmt.Printf("%s(%d) ", next, next.Indent())
		if token.IsEndOfStatementToken(next) {
			fmt.Println()
			printNewLineIndent()
		}
	}
	fmt.Println()
	fmt.Println("Done!")
}
