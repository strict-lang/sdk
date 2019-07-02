package token

import "github.com/BenjaminNitschke/Strict/pkg/source"

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

func (EndOfStatementToken) IsOperator() bool {
	return false
}

func (EndOfStatementToken) IsKeyword() bool {
	return false
}

func (EndOfStatementToken) IsLiteral() bool {
	return false
}

func (EndOfStatementToken) IsValid() bool {
	return true
}

func (EndOfStatementToken) Indent() Indent {
	return 0
}