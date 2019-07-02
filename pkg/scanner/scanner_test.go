package scanner

import (
	"github.com/BenjaminNitschke/Strict/pkg/token"
	"testing"
)

func TestScanner_SkipWhitespaces(test *testing.T) {
	scanner := NewStringScanner("   \t  \r\n  \n")
	scanner.SkipWhitespaces()
	if scanner.lineIndex != 2 {
		test.Errorf("scanner has line-index %d, expected 2", scanner.lineIndex)
	}
}

func TestIndentation(test *testing.T) {
	entries := map[string] token.Indent {
		"  a + b": 2,
		"  b is not true": 2,
		"  b   is   not  true  ": 2,
		"\t\tc is a":4,
		"\t return a":3,
	}
	for entry, indent := range entries {
		scanner := NewStringScanner(entry)
		tokens := scanAllTokens(scanner)
		for _, scanned := range tokens {
			if scanned.Indent() != indent {
				test.Errorf("token %s has indent %d, expected %d", scanned, scanned.Indent(), indent)
			}
		}
	}
}

func TestScanHelloWorld(test *testing.T) {
	const entry = `
	print("Hello, World!");	
	`
	scanner := NewStringScanner(entry)
	assertTokenValue(test, scanner.Pull(), "print")
	assertOperator(test, scanner.Pull(), token.LeftParenOperator)
	assertTokenValue(test, scanner.Pull(), "Hello, World!")
	assertOperator(test, scanner.Pull(), token.RightParenOperator)
	assertOperator(test, scanner.Pull(), token.SemicolonOperator)
	assertEndOfFile(test, scanner.Pull())
}

func TestScanExpression(test *testing.T) {
	const entry = `
	a + b
	`
	scanner := NewStringScanner(entry)
	assertTokenValue(test, scanner.Pull(), "a")
	assertOperator(test, scanner.Pull(), token.AddOperator)
	assertTokenValue(test, scanner.Pull(), "b")
	assertEndOfFile(test, scanner.Pull())
}

func assertTokenValue(test *testing.T, got token.Token, expected string) {
	if got.Value() != expected {
		test.Errorf("unexpected token %s, expected %s", got, expected)
	}
}

func assertOperator(test *testing.T, got token.Token, wanted token.Operator) {
	if !got.IsOperator() || got.(*token.OperatorToken).Operator != wanted {
		test.Errorf("unexpected token %s, expected %s", got, wanted.String())
	}
}

func assertEndOfFile(test *testing.T, got token.Token) {
	_, ok := got.(*token.EndOfFileToken)
	if !ok {
		test.Errorf("unexpected token %s, expected eof", got)
	}
}

func scanAllTokens(scanner *Scanner) []token.Token {
	var tokens []token.Token
	for {
		next := scanner.Pull()
		if _, ok := next.(*token.EndOfFileToken); ok {
			break
		}
		tokens = append(tokens, next)
	}
	return tokens
}
