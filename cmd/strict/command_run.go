package main

import (
	"github.com/spf13/cobra"
)

var runCommand = &cobra.Command{
	Use:   "run",
	Short: "Runs a Strict package",
	Long:  `Runs a package`,
	RunE:  RunRun,
}

var runOptions struct {
	outputPath   string
	reportFormat string
	backendName  string
	debug        bool
}

func init() {
	flags := runCommand.Flags()
	flags.StringVarP(&runOptions.backendName, "backend", "b", "c++", "backend used in code generation")
	flags.BoolVarP(&runOptions.debug, "debug", "z", false, "enable debug mode")
}

func RunRun(command *cobra.Command, arguments []string) error {
	disableLogging()
	fixOptions()
	compilationReport, lineMaps, err := runCompilation()
	if err != nil {
		return err
	}
	output := createOutput(compilationReport, lineMaps)
	return output.Print(command.OutOrStdout())
}
