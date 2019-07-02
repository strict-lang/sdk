package parser

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/token"
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
		return fmt.Sprintf("expected %s but got %s", err.Expected, err.Token)
	}
	return fmt.Sprintf("unexpected token: %s", err.Token)
}

type InvalidIndentationError struct {
	Token token.Token
	Expected token.Indent
}

func (err *InvalidIndentationError) Error() string {
	return fmt.Sprintf(
		"token %s has an invalid indentation level, expected %s",
		err.Token, err.Expected)
}

