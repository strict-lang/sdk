package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.com/strict-lang/sdk/compilation"
	"os"
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
	compiler := &compilation.Compilation{
		Name:          unitName,
		Source:        &compilation.FileSource{File: file},
		TargetArduino: targetArduino,
	}
	result := compiler.Compile()
	if result.Error != nil {
		command.PrintErrf("Failed to compile the file: %s\n", result.Error)
		return
	}
	result.Diagnostics.PrintEntries(&cobraDiagnosticPrinter{command: command})
	if err = writeGeneratedSources(result); err != nil {
		command.PrintErrf("Failed to write generated code; %s\n", err.Error())
		return
	}
	command.Printf("Successfully compiled %s!\n", unitName)
}

func writeGeneratedSources(compilation compilation.Result) (err error) {
	for _, generated := range compilation.GeneratedFiles {
		if err := writeGeneratedSourceFile(generated); err != nil {
			return err
		}
	}
	return nil
}

func writeGeneratedSourceFile(generated compilation.Generated) error {
	file, err := targetFile(generated.FileName)
	if err != nil {
		return err
	}
	_, err = file.Write(generated.Bytes)
	return err
}

func targetFile(name string) (*os.File, error) {
	if buildTargetFile != "" {
		return createNewFile(fmt.Sprintf("./%s", buildTargetFile))
	}
	return createNewFile(fmt.Sprintf("./%s", name))
}
