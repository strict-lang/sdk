package compiler

import (
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"testing"
)

const source = `
String Name
Date Date
`

func TestCompilation(testing *testing.T) {
	result := CompileString(`Test`, source)
	result.Diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
	for _, file := range result.GeneratedFiles {
		testing.Log(string(file.Content))
	}
}
