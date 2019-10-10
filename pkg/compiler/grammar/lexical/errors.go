package lexical

import (
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type UnexpectedCharError struct {
	got      input.Char
	expected string
}

func (err UnexpectedCharError) Error() string {
	return fmt.Sprintf("unexpected char '%c', expected '%s'", err.got, err.expected)
}
