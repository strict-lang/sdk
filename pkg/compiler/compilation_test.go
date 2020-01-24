package compiler

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"testing"
)

func TestCompilation(testing *testing.T) {
	result := CompileString(`Test`,
		`
if True
  if True
		let x = Compute()
		let x = Compute()
		let x = Compute()
		let x = Compute()
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
}
