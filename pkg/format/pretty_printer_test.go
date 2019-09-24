package format

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/parsing"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/scanning"
	"testing"
)

func TestPrettyPrinting(test *testing.T) {
	const source = `
method list<number> rangeTo(number number)
	for num from 1 to number do
		yield num

numbers = rangeTo(100)
for num in numbers do
	if num % 5 is 0 or num % 3 is 0 do
		logf("%d ", num)
`
	unit, err := parsing.NewDefaultFactory().
		WithTokenStream(scanning.NewStringScanning(source)).
		NewParser().
		ParseTranslationUnit()

	if err != nil {
		test.Errorf("unexpected error while parsing: %s", err.Error())
		return
	}

	output := NewStringWriter()
	factory := &PrettyPrinterFactory{
		Unit:   unit,
		Writer: output,
		Format: Format{
			TabWidth:        2,
			ImproveBranches: true,
			IndentWriter: &SimpleSpaceIndentWriter{
				SpacesPerLevel: 2,
			},
			LineLengthLimit: 80,
			EndOfLine:       "\n",
		},
	}
	factory.NewPrettyPrinter().Print()

	test.Logf("Prettified version: \n%s", output.String())
}
