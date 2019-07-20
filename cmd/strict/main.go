package main

import (
	"github.com/spf13/cobra"
	"os"
)

var baseCommand = &cobra.Command{
	Use: "strict",
	Short: "Strict is a CLI for the strict development kit",
	Long: ``,
}

func main() {
	if err := baseCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	baseCommand.AddCommand(buildCommand)
	baseCommand.AddCommand(formatCommand)
}
