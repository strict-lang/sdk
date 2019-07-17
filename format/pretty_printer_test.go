package format

import (
	"gitlab.com/strict-lang/sdk/compiler/parser"
	"gitlab.com/strict-lang/sdk/compiler/scanner"
	"testing"
)

func TestPrettyPrinting(test *testing.T) {
	const source = `
if max(max(1, 2), min(rotate(1), 2)) > 10 do
  print(getSystemOutputStream(), messages.forKey("test.message", getSystemLanguage()))
  return 3213213
else 
	print("NO")
`
	unit, err := parser.
		NewDefaultFactory().
		WithTokenReader(scanner.NewStringScanner(source)).
		NewParser().
		ParseTranslationUnit()

	if err != nil {
		test.Errorf("unexpected error while parsing: %s", err.Error())
		return
	}

	output := NewStringWriter()
	factory := &PrettyPrinterFactory{
		Unit: unit,
		Writer: output,
		Format: Format{
			TabWidth: 2,
			ImproveBranches: true,
			IndentWriter: &SimpleSpaceIndentWriter{
				SpacesPerLevel: 2,
			},
			LineLengthLimit: 80,
			EndOfLine: "\n",
		},
	}
	factory.NewPrettyPrinter().Print()

	test.Logf("Prettified version: \n%s", output.String())
}