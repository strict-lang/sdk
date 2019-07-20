package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/compiler"
)

var buildCommand = &cobra.Command{
	Use: "build [-t target] [-c] [compile flags] [file]",
	Short: "Builds a strict module",
	Long: `Build compiles a file to a specified output file.`,
	Run: RunCompile,
}

var (
	buildOutputFile	string
	compileToCpp bool
)

func init() {
	buildCommand.Flags().
		StringVarP(&buildOutputFile, "target", "t", "", "path to the output file")

	expectNoError(buildCommand.MarkFlagFilename("target", "strict"))
	buildCommand.Flags().BoolVar(&compileToCpp, "c", false, "compile the generated cpp code")
}

func RunCompile(command *cobra.Command, arguments []string) {
	file, ok := findSourceFileInArguments(command, arguments)
	if !ok {
		return
	}
	defer file.Close()
	compilation := compiler.CompileFile(file)
	if compilation.Error != nil {
		command.PrintErrf("Failed to compile the file: %s\n", compilation.Error)
		return
	}
	compilation.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{
		command: command,
	})
	command.Println("Successfully compiled the file")
}


