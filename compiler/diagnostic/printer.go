package diagnostic

import "fmt"

type Printer interface {
	PrintFormatted(format string, arguments ...interface{})
	PrintLine(message string)
	Print(message string)
}

type fmtPrinter struct{}

func NewFmtPrinter() Printer {
	return &fmtPrinter{}
}

func (printer fmtPrinter) Print(message string) {
	fmt.Print(message)
}

func (printer fmtPrinter) PrintLine(message string) {
	fmt.Println(message)
}

func (printer fmtPrinter) PrintFormatted(format string, arguments ...interface{}) {
	fmt.Printf(format, arguments...)
}
