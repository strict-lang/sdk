package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func expectNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func findWorkingDirectory() string {
	directory, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	return directory
}

func findSourceFileInArguments(command *cobra.Command, arguments []string) (*os.File, bool) {
	if len(arguments) == 0 {
		command.PrintErrf("No file given\n")
		return nil, false
	}
	filePath := arguments[0]
	file, err := os.OpenFile(fmt.Sprintf("./%s", filePath), os.O_RDWR, 0)
	if err != nil {
		command.PrintErrf("Failed to open file %s: %s", filePath, err.Error())
		return nil, false
	}
	return file, true
}
