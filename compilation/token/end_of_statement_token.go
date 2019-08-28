package token

import "gitlab.com/strict-lang/sdk/compilation/source"

const (
	EndOfStatementTokenName  = "EndOfStatement"
	EndOfStatementTokenValue = ";"
)

type EndOfStatementToken struct {
	position Position
}

func NewEndOfStatementToken(offset source.Offset) Token {
	return &EndOfStatementToken{
		position: Position{offset, offset},
	}
}

func (token EndOfStatementToken) Position() Position {
	return token.position
}

func (EndOfStatementToken) Name() string {
	return EndOfStatementTokenName
}

func (EndOfStatementToken) Value() string {
	return EndOfStatementTokenValue
}

func (EndOfStatementToken) Indent() Indent {
	return 0
}

func (EndOfStatementToken) String() string {
	return EndOfStatementTokenName
}

func IsEndOfStatementToken(token Token) bool {
	_, ok := token.(*EndOfStatementToken)
	return ok
}
