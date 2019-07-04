package diagnostic

import (
	"fmt"
	"testing"
)

type Printer interface {
	PrintFormatted(format string, arguments ...interface{})
	PrintLine(message string)
	Print(message string)
}

type fmtPrinter struct{}

func NewFmtPrinter() Printer {
	return &fmtPrinter{}
}

func (printer *fmtPrinter) Print(message string) {
	fmt.Print(message)
}

func (printer *fmtPrinter) PrintLine(message string) {
	fmt.Println(message)
}

func (printer *fmtPrinter) PrintFormatted(format string, arguments ...interface{}) {
	fmt.Printf(format, arguments...)
}

type testPrinter struct {
	test *testing.T
}

func NewTestPrinter(test *testing.T) Printer {
	return &testPrinter{
		test: test,
	}
}

func (printer *testPrinter) Print(message string) {
	printer.test.Log(message)
}

func (printer *testPrinter) PrintLine(message string) {
	printer.test.Logf("%s\n", message)
}

func (printer *testPrinter) PrintFormatted(format string, arguments ...interface{}) {
	printer.test.Logf(format, arguments...)
}
