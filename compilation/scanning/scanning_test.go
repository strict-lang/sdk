package scanning

import (
	"gitlab.com/strict-lang/sdk/compilation/token"
	"strings"
	"testing"
)

func TestScannerEndOfStatementInsertion(test *testing.T) {
	entries := map[string]int{
		"add(\na,\nb\n)":                           1,
		"add(a, b)\nadd(b, c)\nadd(b, add(a, c));": 3,
	}

	for entry, expectedCount := range entries {
		scanner := NewStringScanning(entry)
		tokens := ScanAllTokens(scanner)
		count := countEndOfStatements(tokens)
		if count != expectedCount {
			test.Errorf("%s has %d EndOfFile tokens, expected %d", entry, count, expectedCount)
			test.Logf("The tokens: %s", tokens)
		}
	}
}

func countEndOfStatements(tokens []token.Token) int {
	var count int
	for _, element := range tokens {
		if _, ok := element.(*token.EndOfStatementToken); ok {
			count++
		}
	}
	return count
}

func TestIndentation(test *testing.T) {
	entries := map[string]token.Indent{
		"  a + b":                2,
		"  b is not true":        2,
		"  b   is   not  true  ": 2,
	}
	for entry, indent := range entries {
		scanner := NewStringScanning(entry)
		tokens := ScanAllTokens(scanner)
		for _, scanned := range tokens {
			if scanned.Indent() != indent {
				if token.IsEndOfStatementToken(scanned) {
					continue
				}
				test.Errorf("token %s has indent %d, expected %d", scanned, scanned.Indent(), indent)
			}
		}
	}
}

func TestScanHelloWorld(test *testing.T) {
	const entry = `
	print(string(1))
	`
	scanner := NewStringScanning(entry)
	assertTokenValue(test, scanner.Pull(), "print")
	assertOperator(test, scanner.Pull(), token.LeftParenOperator)
	assertTokenValue(test, scanner.Pull(), "string")
	assertOperator(test, scanner.Pull(), token.LeftParenOperator)
	assertTokenValue(test, scanner.Pull(), "1")
	assertOperator(test, scanner.Pull(), token.RightParenOperator)
	assertOperator(test, scanner.Pull(), token.RightParenOperator)
	assertEndOfStatement(test, scanner.Pull())
	assertEndOfFile(test, scanner.Pull())
}

func TestScanExpression(test *testing.T) {
	const entry = `
	a + b
	`
	scanner := NewStringScanning(entry)
	assertTokenValue(test, scanner.Pull(), "a")
	assertOperator(test, scanner.Pull(), token.AddOperator)
	assertTokenValue(test, scanner.Pull(), "b")
	assertEndOfStatement(test, scanner.Pull())
	assertEndOfFile(test, scanner.Pull())
}

func TestScanMethodDeclaration(test *testing.T) {
	const entry = `
	method number add(number a, number b) 
		return 0.123 + 3.210
	`
	scanner := NewStringScanning(entry)
	assertKeyword(test, scanner.Pull(), token.MethodKeyword)
	assertTokenValue(test, scanner.Pull(), "number")
	assertTokenValue(test, scanner.Pull(), "add")
	assertOperator(test, scanner.Pull(), token.LeftParenOperator)
	assertTokenValue(test, scanner.Pull(), "number")
	assertTokenValue(test, scanner.Pull(), "a")
	assertOperator(test, scanner.Pull(), token.CommaOperator)
	assertTokenValue(test, scanner.Pull(), "number")
	assertTokenValue(test, scanner.Pull(), "b")
	assertOperator(test, scanner.Pull(), token.RightParenOperator)
	assertEndOfStatement(test, scanner.Pull())
	assertKeyword(test, scanner.Pull(), token.ReturnKeyword)
	assertTokenValue(test, scanner.Pull(), "0.123")
	assertOperator(test, scanner.Pull(), token.AddOperator)
	assertTokenValue(test, scanner.Pull(), "3.210")
	assertEndOfStatement(test, scanner.Pull())
	assertEndOfFile(test, scanner.Pull())
}

func assertTokenValue(test *testing.T, got token.Token, expected string) {
	if got.Value() != expected {
		test.Errorf("unexpected token %s, expected %s", got, expected)
	}
}

func assertOperator(test *testing.T, got token.Token, wanted token.Operator) {
	if token.OperatorValue(got) != wanted {
		test.Errorf("unexpected token %s, expected %s", got, wanted.String())
	}
}

func assertKeyword(test *testing.T, got token.Token, wanted token.Keyword) {
	if token.KeywordValue(got) != wanted {
		test.Errorf("unexpected token %s, expected %s", got, wanted.String())
	}
}

func assertEndOfStatement(test *testing.T, got token.Token) {
	if !token.IsEndOfStatementToken(got) {
		test.Errorf("unexpected token %s, expected %s", got, token.EndOfStatementTokenName)
	}
}

func assertEndOfFile(test *testing.T, got token.Token) {
	if !token.IsEndOfFileToken(got) {
		test.Errorf("unexpected token %s, expected %s", got, token.EndOfFileTokenName)
	}
}

func sumToIndex(array []int, target int) (sum int) {
	for index := 0; index < target; index++ {
		sum += array[index]
	}
	return
}

func createTextWithLineLengths(lengths []int) string {
	var builder strings.Builder
	for _, line := range lengths {
		for count := 0; count < line; count++ {
			builder.WriteRune(' ')
		}
		builder.WriteRune('\n')
	}
	return builder.String()
}
