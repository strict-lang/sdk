package main

import (
	"github.com/spf13/cobra"
)

type cobraDiagnosticPrinter struct {
	command *cobra.Command
}

func (printer *cobraDiagnosticPrinter) Print(message string) {
	printer.command.Print(message)
}

func (printer *cobraDiagnosticPrinter) PrintLine(message string) {
	printer.command.Println(message)
}

func (printer *cobraDiagnosticPrinter) PrintFormatted(format string, arguments ...interface{}) {
	printer.command.Printf(format, arguments...)
}
