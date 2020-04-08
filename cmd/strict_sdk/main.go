package main

import (
	"github.com/spf13/cobra"
	"os"
)

var baseCommand = &cobra.Command{
	Use:   "strict_sdk",
	Short: "The SDK command is used to build a Strict strict_sdk",
}

func main() {
	if err := baseCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	baseCommand.AddCommand(makeCommand)
}
