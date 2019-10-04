package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/pkg/compilation/scanning"
	"gitlab.com/strict-lang/sdk/pkg/compilation/source"
	"gitlab.com/strict-lang/sdk/pkg/compilation/token"
	"os"
)

var tokenizeCommand = &cobra.Command{
	Use:   "tokenize [-c] [file]",
	Short: "Tokenizes the file and prints the tokens",
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
	scan := scanning.NewScanning(source.NewStreamReader(sourceFile))
	fmt.Println("Scanned tokens:")
	printNewLineIndent()
	for {
		next := scan.Pull()
		if token.IsEndOfFileToken(next) {
			break
		}
		fmt.Printf("%s ", next)
		if token.IsEndOfStatementToken(next) {
			fmt.Println()
			printNewLineIndent()
		}
	}
	fmt.Println()
	fmt.Println("Done!")
}