package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/diagnostic"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"github.com/BenjaminNitschke/Strict/compiler/token"
	"testing"
)

func NewTestParser(tokens token.Reader) *Parser {
	return NewParser("test", tokens, diagnostic.NewRecorder())
}

func TestParseTopLevelStatements(test *testing.T) {
	const entry = `
method list<number> divisibleNumbers(number limit)
  for index from 0 to limit do
		if index % 3 is 0 or index % 5 is 0
			yield index

numbers = divisibleNumbers(10)
for element in numbers do
  logFormatted("%d", toInt(element))
`
	parser := NewTestParser(scanner.NewStringScanner(entry))
	parser.parseTopLevelNodes()
}
