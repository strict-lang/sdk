package token

type Token interface {
	Name() string
	Value() string
	Position() Position
	IsOperator() bool
	IsKeyword() bool
	IsLiteral() bool
	IsValid() bool
}