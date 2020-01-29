package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"strict.dev/sdk/pkg/compiler"
	"os"
	"strict.dev/sdk/pkg/compiler/backend"
)

var buildCommand = &cobra.Command{
	Use:   "build [-t target] [-a arduino] [-c] [compile flags] [file]",
	Short: "Builds a Strict module",
	Long:  `Build compiles a file to a specified output file.`,
	Run:   RunCompile,
}

var (
	buildTargetFile string
	targetArduino   bool
	compileToCpp    bool
)

func init() {
	buildCommand.Flags().
		StringVarP(&buildTargetFile, "target", "t", "", "path to the output file")

	expectNoError(buildCommand.MarkFlagFilename("target", "Strict"))
	buildCommand.Flags().BoolVarP(&targetArduino, "arduino", "a", false, "generate arduino code")
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
	compiler := &compiler.Compilation{
		Name:          unitName,
		Source:        &compiler.FileSource{File: file},
		TargetArduino: targetArduino,
	}
	result := compiler.Compile()
	result.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{command: command})
	if result.Error != nil {
		command.PrintErrf("The compilation has failed: %s\n", result.Error)
		return
	}
	if err = writeGeneratedSources(result); err != nil {
		command.PrintErrf("Failed to write generated code; %s\n", err.Error())
		return
	}
	command.Printf("Successfully compiled %s!\n", unitName)
}

func writeGeneratedSources(compilation compiler.Result) (err error) {
	for _, generated := range compilation.GeneratedFiles {
		if err = writeGeneratedSourceFile(generated); err != nil {
			return err
		}
	}
	return nil
}

func writeGeneratedSourceFile(generated backend.GeneratedFile) error {
	file, err := targetFile(generated.Name)
	if err != nil {
		return err
	}
	_, err = file.Write(generated.Content)
	return err
}

func targetFile(name string) (*os.File, error) {
	if buildTargetFile != "" {
		return createNewFile(fmt.Sprintf("./%s", buildTargetFile))
	}
	return createNewFile(fmt.Sprintf("./%s", name))
}
