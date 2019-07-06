package parser

import (
	"testing"

	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"github.com/BenjaminNitschke/Strict/compiler/token"
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
