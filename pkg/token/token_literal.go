package token

const (
	StringLiteralTokenName = "StringLiteral"
	NumberLiteralTokenName = "NumberLiteral"
)

type LiteralToken struct {
	name     string
	value    string
	position Position
}

func NewStringLiteralToken(value string, position Position) LiteralToken {
	return LiteralToken{
		name:     StringLiteralTokenName,
		value:    value,
		position: position,
	}
}

func NewNumberLiteralToken(value string, position Position) LiteralToken {
	return LiteralToken{
		name:     NumberLiteralTokenName,
		value:    value,
		position: position,
	}
}

func (literal LiteralToken) Name() string {
	return literal.name
}

func (literal LiteralToken) Value() string {
	return literal.value
}

func (literal LiteralToken) Position() Position {
	return literal.position
}

func (LiteralToken) IsKeyword() bool {
	return false
}

func (LiteralToken) IsOperator() bool {
	return false
}

func (LiteralToken) IsLiteral() bool {
	return true
}

func (LiteralToken) IsValid() bool {
	return true
}
