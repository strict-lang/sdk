package token

type Operator string

const (
	AddOperator Operator = "+"
)

const	OperatorTokenName = "operator"

type OperatorToken struct {
	operator Operator
	position Position
}

func NewOperatorToken(operator Operator, position Position) Token {
	return &OperatorToken{
		operator: operator,
		position: position,
	}
}

func (operator OperatorToken) Value() string {
	return string(operator.operator)
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