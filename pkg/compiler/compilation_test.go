package compiler

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"testing"
)

func TestCompilation(testing *testing.T) {
	result := CompileString(`Test`,
		`
method Compute()
	if True
    if True
			let x = Compute()
			let x = Compute()
			let x = Compute()
			let x = Compute()
    Compute()
    let y = compute()
    for i from 0 to 1
      Compute(i)
	if True
    if True
			let x = Compute()
	let x = Compute()
	let x = Compute()
	let x = Compute()
	let x = Compute()
	let x = Compute()
`)
	result.Diagnostics.PrintEntries(diagnostic.NewFmtPrinter())
	for _, file := range result.GeneratedFiles {
		fmt.Println(string(file.Bytes))
	}
}
