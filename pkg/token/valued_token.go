package token

const (
	StringLiteralTokenName = "StringLiteralToken"
	NumberLiteralTokenName = "NumberLiteralToken"
	IdentifierTokenName    = "IdentifierToken"
)

type ValuedToken struct {
	name     string
	value    string
	position Position
	literal  bool
}

func NewStringLiteralToken(value string, position Position) *ValuedToken {
	return &ValuedToken{
		name:     StringLiteralTokenName,
		value:    value,
		position: position,
		literal:  true,
	}
}

func NewNumberLiteralToken(value string, position Position) *ValuedToken {
	return &ValuedToken{
		name:     NumberLiteralTokenName,
		value:    value,
		position: position,
		literal:  true,
	}
}

func NewIdentifierToken(value string, position Position) *ValuedToken {
	return &ValuedToken{
		name:     IdentifierTokenName,
		value:    value,
		position: position,
		literal:  true,
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

func (ValuedToken) IsKeyword() bool {
	return false
}

func (ValuedToken) IsOperator() bool {
	return false
}

func (token ValuedToken) IsLiteral() bool {
	return token.literal
}

func (ValuedToken) IsValid() bool {
	return true
}
