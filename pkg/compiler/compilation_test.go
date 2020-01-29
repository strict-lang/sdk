package compiler

import (
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"testing"
)

const source = `
method Number Compute()
	if True
		let y = Compute()
		for index from 0 to 1
			Compute(index)
	if let z = 10 + 20
		if True
			let x = Compute()
`

func TestCompilation(testing *testing.T) {
	result := CompileString(`Test`, source)
	result.Diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
	for _, file := range result.GeneratedFiles {
		testing.Log(string(file.Content))
	}
}
