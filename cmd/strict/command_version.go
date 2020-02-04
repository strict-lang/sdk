package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Version = "0.2.0"

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the installed SDKs version",
	Run:   PrintVersion,
}

func PrintVersion(command *cobra.Command, arguments []string) {
	fmt.Println("Strict version ", Version)
}
