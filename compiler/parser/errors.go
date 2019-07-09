package parser

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

// UnexpectedTokenError indicates that the parser expected a certain kind of token, but
// got a different one. It captures the token and has an optional 'expected' field, which
// stores the name of the kind of token that was expected.
type UnexpectedTokenError struct {
	Token    token.Token
	Expected string
}

func (err *UnexpectedTokenError) Error() string {
	if err.Expected != "" {
		return fmt.Sprintf("expected %s but got {%s}", err.Expected, err.Token)
	}
	return fmt.Sprintf("unexpected token: {%s}", err.Token)
}

// InvalidIndentationError indicates that the indentation of a token
// in a block of statements is invalid. The tokens indent always has
// to match that of its block.
type InvalidIndentationError struct {
	Token    token.Token
	Expected string
}

func (err *InvalidIndentationError) Error() string {
	return fmt.Sprintf(
		"token {%s} has an invalid indentation level of %d, expected %s",
		err.Token, err.Token.Indent(), err.Expected)
}

// InvalidStatementError indicates that the statement in the line is invalid
// and could not be parsed.
type InvalidStatementError struct {
	LineIndex source.LineIndex
}

func (err *InvalidStatementError) Error() string {
	return fmt.Sprintf(
		"invalid statement in line %d", err.LineIndex)
}
