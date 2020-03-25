package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"gitlab.com/strict-lang/sdk/pkg/compiler"
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
)

var buildCommand = &cobra.Command{
	Use:   "build [",
	Short: "Builds a Strict module",
	Long:  `Build compiles a file to a specified output file.`,
	Run:   RunCompile,
}

var buildOptions struct {
	backendId  string
	outputPath string
}

func init() {
	flags := buildCommand.Flags()
	flags.StringVarP(&buildOptions.outputPath, "output", "o", "", "output path")
	flags.StringVarP(&buildOptions.backendId, "backend", "b", "cpp", "id of the backend")
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
	compilation := &compiler.Compilation{
		Name:    unitName,
		Source:  &compiler.FileSource{File: file},
		Backend: buildOptions.backendId,
	}
	result := compilation.Compile()
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
	if buildOptions.outputPath != "" {
		return createNewFile(fmt.Sprintf("./%s", buildOptions.outputPath))
	}
	return createNewFile(fmt.Sprintf("./%s", name))
}
