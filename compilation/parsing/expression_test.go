package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/scanning"
	"testing"
)

func TestParseBinaryExpression(test *testing.T) {
	entries := []string{
		"call(call(arg))",
		"1 isnt 1",
		"1 is 1 or 1 isnt 2",
		"random % 2 is 1",
		"index % 3 is 0 or index % 5 is 0",
		"!1",
		"(1 + 2)",
		"1 + 2 + 3",
		"(1 + 2) + 3",
		"printf(\"%d\", limit(10) + 1))",
	}

	for _, entry := range entries {
		testParsingBinaryExpression(test, entry)
	}
}

func testParsingBinaryExpression(test *testing.T, entry string) {
	parser := NewTestParser(scanning.NewStringScanning(entry))
	_, err := parser.parseExpression()
	if err != nil {
		test.Errorf("unexpected error while parsing (%s): %s", entry, err.Error())
		return
	}
}
