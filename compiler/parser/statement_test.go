package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/scanner"
	"github.com/BenjaminNitschke/Strict/compiler/token"
	"testing"
)

func TestParser_ParseMethodCall(test *testing.T) {
	const entry = `
		add(a, b)
	`

	tokens := scanner.NewStringScanner(entry)
	parser := NewTestParser(tokens)
	// The first token has to be pulled, ParseMethodCall will start
	// with the 'last' token. It is normally called after the identifier
	// has already been pulled, because one has to peek one token to
	// check whether an identifier is actually belonging to a method-call
	// and the scanner does not allow peeking further than one token.
	parser.tokens.Pull()
	call, err := parser.ParseMethodCall()
	if err != nil {
		test.Errorf("unexpected error: %s", err)
		return
	}
	if call.Name.Value != "add" {
		test.Errorf("methodCall has name %s, expected 'add'", call.Name.Value)
	}
	if arguments := len(call.Arguments); arguments != 2 {
		test.Errorf("methodCall has %d argument(s), expected 2", arguments)
	}
	endOfStatement := parser.tokens.Pull()
	if !token.IsEndOfStatementToken(endOfStatement) {
		test.Errorf(
			"unexpected token %s, expected %s",
			endOfStatement, token.EndOfStatementTokenName)
	}
	endOfFile := parser.tokens.Pull()
	if !token.IsEndOfFileToken(endOfFile) {
		test.Errorf(
			"unexpected token %s, expected %s",
			endOfStatement, token.EndOfStatementTokenName)
	}
}
