package parsing

import (
	"testing"

	"gitlab.com/strict-lang/sdk/compilation/scanning"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

func NewTestParser(tokens token.Stream) *Parsing {
	return NewDefaultFactory().WithTokenStream(tokens).NewParser()
}

func TestParseTopLevelStatements(test *testing.T) {
	const entry = `
import "stdio.h" as Io
import "something"

shared int a = 0
shared int x

method nothing()

method list<number> range(number begin, number end)
  for num from begin to end do
    yield num

for num in range(1, 21) do
  if num % 3 is 0 and num % 5 is 0 do
    stdio.puts("FizzBuzz")
  else if num % 3 is 0 do
		stdio.puts("Fizz")
	else if num % 5 is 0 do
		stdio.puts("Buzz")
	else
		stdio.printf("%d", num)
`
	parser := NewTestParser(scanning.NewStringScanning(entry))
	parser.parseTopLevelNodes()
}
