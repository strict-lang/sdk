package token

import "fmt"

type Token struct {
	Value    string
	Kind     Kind
	Position Position
}

var InvalidToken = Token{
	Value:    "",
	Kind:     Invalid,
	Position: Position{},
}

func (token Token) String() string {
	group := token.Kind.Group()
	return fmt.Sprintf("%s{%s}", group, token.Value)
}

func NewStringLiteral(text string, position Position) Token {
	return Token{
		Value:    text,
		Position: position,
		Kind:     StringLiteral,
	}
}

func NewOperatorToken(operator Kind, position Position) Token {
	return Token{
		Kind:     operator,
		Value:    NameOfKind(operator),
		Position: position,
	}
}

func NewKeywordToken(keyword Kind, position Position) Token {
	return Token{
		Kind:     keyword,
		Value:    NameOfKind(keyword),
		Position: position,
	}
}
