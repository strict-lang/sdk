package token

const InvalidTokenName = "invalid"

type InvalidToken struct {
	value    string
	position Position
}

func NewInvalidToken(value string, position Position) Token {
	return &InvalidToken{
		value: value,
		position: position,
	}
}

func (invalid InvalidToken) Name() string {
	return InvalidTokenName
}

func (invalid InvalidToken) Value() string {
	return invalid.value
}

func (invalid InvalidToken) Position() Position {
	return invalid.Position()
}

func (InvalidToken) IsValid() bool {
	return false
}

func (InvalidToken) IsKeyword() bool {
	return false
}

func (InvalidToken) IsOperator() bool {
	return false
}

func (InvalidToken) IsLiteral() bool {
	return false
}