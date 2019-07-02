package token

type Indent uint8

type Token interface {
	Name() string
	Value() string
	Position() Position
	Indent() Indent
}
