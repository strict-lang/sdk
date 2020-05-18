package main

import (
	"github.com/spf13/cobra"
	"os"
)

var baseCommand = &cobra.Command{
	Use:   "Strict",
	Short: "Strict is a CLI for the Strict development kit",
	Long:  ``,
}

func main() {
	if err := baseCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	baseCommand.AddCommand(buildCommand)
	baseCommand.AddCommand(versionCommand)
	baseCommand.AddCommand(treeCommand)
	baseCommand.AddCommand(tokenizeCommand)
	baseCommand.AddCommand(initCommand)
	baseCommand.AddCommand(runCommand)
}
