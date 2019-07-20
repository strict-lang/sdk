package main

import (
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/compiler"
	"os"
)

var buildCommand = &cobra.Command{
	Use:   "build [-t target] [-c] [compile flags] [file]",
	Short: "Builds a strict module",
	Long:  `Build compiles a file to a specified output file.`,
	Run:   RunCompile,
}

var (
	buildTargetFile string
	compileToCpp    bool
)

func init() {
	buildCommand.Flags().
		StringVarP(&buildTargetFile, "target", "t", "", "path to the output file")

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
	if err := writeGeneratedSources(compilation); err != nil {
		command.PrintErrf("Failed to write generated code; %s", err.Error())
		return
	}
	command.Println("Successfully compiled the file")
}

func writeGeneratedSources(compilation compiler.CompilationResult) (err error) {
	file, err := targetFile(compilation.UnitName)
	_, err = file.Write(compilation.Generated)
	return nil
}

func targetFile(unitName string) (*os.File, error) {
	if buildTargetFile != "" {
		return os.Open(buildTargetFile)
	}
	name := GeneratedFileName(unitName)
	return os.Open(name)
}
