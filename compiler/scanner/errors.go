package scanner

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/source"
)

type UnexpectedCharError struct {
	got      source.Char
	expected string
}

func (err UnexpectedCharError) Error() string {
	return fmt.Sprintf("unexpected char '%c', expected '%s'", err.got, err.expected)
}
