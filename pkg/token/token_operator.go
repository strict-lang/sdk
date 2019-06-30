package token

import "fmt"

type Operator int8

const (
	AddOperator Operator = iota
	SubOperator
	MulOperator
	DivOperator
	ModOperator
	EqualsOperator
	NotEqualsOperator
	ShiftLeftOperator
	ShiftRightOperator
	AndOperator
	XorOperator
	OrOperator
	GreaterOperator
	GreaterEqualsOperator
	AssignOperator
	ColonOperator
	SmallerOperator
	SmallerEqualsOperator
)

const OperatorTokenName = "operator"

type OperatorToken struct {
	Operator Operator
	position Position
}

func NewOperatorToken(operator Operator, position Position) Token {
	return &OperatorToken{
		Operator: operator,
		position: position,
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

type Precedence int8

const (
	LowPrecedence   = 0
	UnaryPrecedence = 6
	HighPrecedence  = 7
)

func (operator Operator) Precedence() int {
	switch operator {
	case EqualsOperator,
		NotEqualsOperator,
		GreaterOperator,
		GreaterEqualsOperator:
		return 3
	case AddOperator,
		SubOperator,
		OrOperator,
		XorOperator:
		return 4
	case MulOperator,
		DivOperator,
		ModOperator,
		ShiftLeftOperator,
		ShiftRightOperator,
		AndOperator:
		return LowPrecedence
	}
	return 0
}

func (operator Operator) String() string {
	return fmt.Sprintf("operator(%d)", operator)
}