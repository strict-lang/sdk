package syntax

import (
	"errors"
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/grammar/token"
)

var (
	InvalidStatementError = errors.New("invalid statement")
	// ErrInvalidExpression is returned from a function that fails to grammar
	// an expression. Functions returning this should report more verbose
	// error messages to the diagnostics.Bag.
	ErrInvalidExpression = errors.New("invalid expression")
)

// UnexpectedTokenError indicates that the grammar expected a certain kind of token, but
// got a different one. It captures the token and has an optional 'expected' field, which
// stores the name of the kind of token that was expected.
type UnexpectedTokenError struct {
	Token    token.Token
	Expected string
}

func (err *UnexpectedTokenError) Error() string {
	if err.Expected != "" {
		return fmt.Sprintf("expected %s but got: '%s'", err.Expected, err.Token)
	}
	return fmt.Sprintf("unexpected token: '%s'", err.Token)
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
		"invalid indent of %d, expected %s",
		err.Token.Indent(), err.Expected)
}
