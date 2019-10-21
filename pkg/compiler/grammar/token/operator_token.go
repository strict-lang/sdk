package token

import "fmt"

const OperatorTokenName = "op"

type OperatorToken struct {
	Operator Operator
	position Position
	indent   Indent
}

func NewOperatorToken(operator Operator, position Position, indent Indent) Token {
	return &OperatorToken{
		Operator: operator,
		position: position,
		indent:   indent,
	}
}

func (operator OperatorToken) Value() string {
	return string(operator.Operator)
}

func (operator OperatorToken) Position() Position {
	return operator.position
}

func (OperatorToken) Name() string {
	return OperatorTokenName
}

func (OperatorToken) IsOperator() bool {
	return true
}

func (OperatorToken) IsKeyword() bool {
	return false
}

func (OperatorToken) IsLiteral() bool {
	return false
}

func (OperatorToken) IsValid() bool {
	return true
}

func (operator OperatorToken) Indent() Indent {
	return operator.indent
}

func (operator OperatorToken) String() string {
	return fmt.Sprintf("%s(%s)", OperatorTokenName, operator.Operator)
}

func IsOperatorToken(token Token) bool {
	_, ok := token.(*OperatorToken)
	return ok
}

func PrecedenceOfAny(token Token) Precedence {
	// Will be LowestPrecedence if the token has no operator value
	return OperatorValue(token).Precedence()
}
