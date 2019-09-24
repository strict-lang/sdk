package scanning

import (
	"fmt"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

type UnexpectedCharError struct {
	got      source2.Char
	expected string
}

func (err UnexpectedCharError) Error() string {
	return fmt.Sprintf("unexpected char '%c', expected '%s'", err.got, err.expected)
}
