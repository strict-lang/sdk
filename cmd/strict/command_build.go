package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/compiler"
	"os"
)

var buildCommand = &cobra.Command{
	Use:   "build [-t target] [-p platform] [-c] [compile flags] [file]",
	Short: "Builds a strict module",
	Long:  `Build compiles a file to a specified output file.`,
	Run:   RunCompile,
}

var (
	buildTargetFile string
	targetPlatform  string
	compileToCpp    bool
)

func init() {
	buildCommand.Flags().
		StringVarP(&buildTargetFile, "target", "t", "", "path to the output file")

	buildCommand.Flags().
		StringVarP(&targetPlatform, "platform", "p", "cross", "name of the target platform")
	expectNoError(buildCommand.MarkFlagFilename("target", "strict"))
	buildCommand.Flags().BoolVar(&compileToCpp, "c", false, "compile the generated cpp code")
}

func RunCompile(command *cobra.Command, arguments []string) {
	file, ok := findSourceFileInArguments(command, arguments)
	if !ok {
		return
	}
	defer file.Close()
	unitName, err := ParseUnitName(file.Name())
	if err != nil {
		command.Printf("Invalid filename: %s\n", file.Name())
		return
	}
	compilation := compiler.CompileFile(unitName, file)
	if compilation.Error != nil {
		command.PrintErrf("Failed to compile the file: %s\n", compilation.Error)
		return
	}
	compilation.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{
		command: command,
	})
	if err := writeGeneratedSources(compilation); err != nil {
		command.PrintErrf("Failed to write generated code; %s\n", err.Error())
		return
	}
	command.Printf("Successfully compiled %s!\n", unitName)
}

func writeGeneratedSources(compilation compiler.CompilationResult) (err error) {
	file, err := targetFile(compilation.UnitName)
	if err != nil {
		return
	}
	_, err = file.Write(compilation.Generated)
	return
}

func targetFile(unitName string) (*os.File, error) {
	if buildTargetFile != "" {
		return createNewFile(fmt.Sprintf("./%s", buildTargetFile))
	}
	name := GeneratedFileName(unitName)
	return createNewFile(fmt.Sprintf("./%s",name))
}
