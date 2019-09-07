package token

import "fmt"

const InvalidTokenName = "invalid"

type InvalidToken struct {
	value    string
	position Position
	indent   Indent
}

func NewAnonymousInvalidToken() Token {
	return &InvalidToken{value: "", indent: 0}
}

func NewInvalidToken(value string, position Position, indent Indent) Token {
	return &InvalidToken{
		value:    value,
		position: position,
		indent:   indent,
	}
}

func (invalid InvalidToken) Name() string {
	return InvalidTokenName
}

func (invalid InvalidToken) Value() string {
	return invalid.value
}

func (invalid InvalidToken) Position() Position {
	return invalid.position
}

func (invalid InvalidToken) Indent() Indent {
	return invalid.indent
}

func (invalid InvalidToken) String() string {
	if invalid.value == "" {
		return InvalidTokenName
	}
	return fmt.Sprintf("%s(%s)", InvalidTokenName, invalid.value)
}

func IsInvalidToken(token Token) bool {
	_, ok := token.(*InvalidToken)
	return ok
}
