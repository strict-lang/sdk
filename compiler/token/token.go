package token

type Indent uint8
const NoIndent Indent = 0

type Token interface {
	Name() string
	Value() string
	Position() Position
	Indent() Indent
}
