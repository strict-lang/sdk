package token

type Indent uint8

type Token interface {
	Name() string
	Value() string
	Position() Position
	IsOperator() bool
	IsKeyword() bool
	IsLiteral() bool
	IsValid() bool
	Indent() Indent
}
