package token

type LiteralToken struct {
	name     string
	value    string
	position Position
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
