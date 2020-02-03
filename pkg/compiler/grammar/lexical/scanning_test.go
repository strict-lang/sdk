package lexical

import (
	"strict.dev/sdk/pkg/compiler/grammar/token"
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
		tokens := scanRemaining(scanner)
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
		"  b   is   not  true  ": 2,
		"  a + b":                2,
		"  b is not true":        2,
	}
	for entry, indent := range entries {
		scanner := NewStringScanning(entry)
		tokens := scanRemaining(scanner)
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
	assertTokenValue(test, scanner.Pull(), "print", entry)
	assertOperator(test, scanner.Pull(), token.LeftParenOperator, entry)
	assertTokenValue(test, scanner.Pull(), "string", entry)
	assertOperator(test, scanner.Pull(), token.LeftParenOperator, entry)
	assertTokenValue(test, scanner.Pull(), "1", entry)
	assertOperator(test, scanner.Pull(), token.RightParenOperator, entry)
	assertOperator(test, scanner.Pull(), token.RightParenOperator, entry)
	assertEndOfStatement(test, scanner.Pull(), entry)
	assertEndOfFile(test, scanner.Pull(), entry)
}

func TestScanExpression(test *testing.T) {
	const entry = `
	a + b
	`
	scanner := NewStringScanning(entry)
	assertTokenValue(test, scanner.Pull(), "a", entry)
	assertOperator(test, scanner.Pull(), token.AddOperator, entry)
	assertTokenValue(test, scanner.Pull(), "b", entry)
	assertEndOfStatement(test, scanner.Pull(), entry)
	assertEndOfFile(test, scanner.Pull(), entry)
}

func TestScanMethodDeclaration(test *testing.T) {
	const entry = `
	method number add(number a, number b) 
		return 0.123 + 3.210
	`
	scanner := NewStringScanning(entry)
	assertKeyword(test, scanner.Pull(), token.MethodKeyword, entry)
	assertTokenValue(test, scanner.Pull(), "number", entry)
	assertTokenValue(test, scanner.Pull(), "add", entry)
	assertOperator(test, scanner.Pull(), token.LeftParenOperator, entry)
	assertTokenValue(test, scanner.Pull(), "number", entry)
	assertTokenValue(test, scanner.Pull(), "a", entry)
	assertOperator(test, scanner.Pull(), token.CommaOperator, entry)
	assertTokenValue(test, scanner.Pull(), "number", entry)
	assertTokenValue(test, scanner.Pull(), "b", entry)
	assertOperator(test, scanner.Pull(), token.RightParenOperator, entry)
	assertEndOfStatement(test, scanner.Pull(), entry)
	assertKeyword(test, scanner.Pull(), token.ReturnKeyword, entry)
	assertTokenValue(test, scanner.Pull(), "0.123", entry)
	assertOperator(test, scanner.Pull(), token.AddOperator, entry)
	assertTokenValue(test, scanner.Pull(), "3.210", entry)
	assertEndOfStatement(test, scanner.Pull(), entry)
	assertEndOfFile(test, scanner.Pull(), entry)
}

func assertTokenValue(test *testing.T, got token.Token, expected string, entry string) {
	if got.Value() != expected {
		test.Errorf("unexpected token while scanning \n%s\n %s, expected %s",
			entry, got, expected)
	}
}

func assertOperator(test *testing.T, got token.Token, wanted token.Operator, entry string) {
	if token.OperatorValue(got) != wanted {
		test.Errorf("unexpected token while scanning \n%s\n %s, expected %s",
			entry, got, wanted.String())
	}
}

func assertKeyword(test *testing.T, got token.Token, wanted token.Keyword, entry string) {
	if token.KeywordValue(got) != wanted {
		test.Errorf("unexpected token while scanning \n%s\n expected %s but got %s",
			entry, got, wanted.String())
	}
}

func assertEndOfStatement(test *testing.T, got token.Token, entry string) {
	if !token.IsEndOfStatementToken(got) {
		test.Errorf("unexpected token while scanning \n%s\n expected %s but got %s",
			entry, got, token.EndOfStatementTokenName)
	}
}

func assertEndOfFile(test *testing.T, got token.Token, entry string) {
	if !token.IsEndOfFileToken(got) {
		test.Errorf("unexpected token while scanning \n%s\n expected %s but got %s",
			entry, token.EndOfFileTokenName, got)
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
