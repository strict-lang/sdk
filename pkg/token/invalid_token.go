package token

const InvalidTokenName = "invalid"

type InvalidToken struct {
	value    string
	position Position
	indent Indent
}

func NewAnonymousInvalidToken() Token {
	return &InvalidToken{value: "", indent: 0}
}

func NewInvalidToken(value string, position Position, indent Indent) Token {
	return &InvalidToken{
		value:    value,
		position: position,
		indent: indent,
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

func (invalid InvalidToken) Indent() Indent {
	return invalid.indent
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
