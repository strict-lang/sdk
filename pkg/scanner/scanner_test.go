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
