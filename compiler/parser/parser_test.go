package parser

import (
	"testing"

	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/scanner"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

func NewTestParser(tokens token.Reader) *Parser {
	return NewDefaultFactory().WithTokenReader(tokens).NewParser()
}

func TestParseTopLevelStatements(test *testing.T) {
	const entry = `
method list<number> range(number begin, number end)
  for num from begin to end do
    yield num

for num in range(1, 21) do
  if num % 3 is 0 and num % 5 is 0
    log("FizzBuzz")
  else 
		if num % 3 is 0
    	log("Fizz")
  	else 
			if num % 5 is 0
    		log("Buzz")
			else
    		logf("%d", num)
`
	parser := NewTestParser(scanner.NewStringScanner(entry))
	nodes := parser.parseTopLevelNodes()
	for _, node := range nodes {
		ast.Print(node)
	}
}
