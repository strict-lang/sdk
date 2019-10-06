package token

import (
	"fmt"
)

const (
	StringLiteralTokenName = "str"
	NumberLiteralTokenName = "num"
	IdentifierTokenName    = "id"
)

type ValuedToken struct {
	name     string
	value    string
	position Position
	literal  bool
	indent   Indent
}

func NewStringLiteralToken(value string, position Position, indent Indent) *ValuedToken {
	return &ValuedToken{
		name:     StringLiteralTokenName,
		value:    value,
		position: position,
		literal:  true,
		indent:   indent,
	}
}

func NewNumberLiteralToken(value string, position Position, indent Indent) *ValuedToken {
	return &ValuedToken{
		name:     NumberLiteralTokenName,
		value:    value,
		position: position,
		literal:  true,
		indent:   indent,
	}
}

func NewIdentifierToken(value string, position Position, indent Indent) *ValuedToken {
	return &ValuedToken{
		name:     IdentifierTokenName,
		value:    value,
		position: position,
		literal:  false,
		indent:   indent,
	}
}

func (token ValuedToken) Name() string {
	return token.name
}

func (token ValuedToken) Value() string {
	return token.value
}

func (token ValuedToken) Position() Position {
	return token.position
}

func (token ValuedToken) Indent() Indent {
	return token.indent
}

func (token ValuedToken) String() string {
	return fmt.Sprintf("%s(%s)", token.name, token.value)
}

func IsLiteralToken(token Token) bool {
	valued, ok := token.(*ValuedToken)
	if !ok {
		return false
	}
	return valued.literal
}

func IsStringLiteralToken(token Token) bool {
	return token.Name() == StringLiteralTokenName
}

func IsNumberLiteralToken(token Token) bool {
	return token.Name() == NumberLiteralTokenName
}

func IsIdentifierToken(token Token) bool {
	valued, ok := token.(*ValuedToken)
	if !ok {
		return false
	}
	return !valued.literal
}
