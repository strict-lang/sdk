package token

type Operator int8

const (
	InvalidOperator Operator = iota
	unaryOperatorBegin
	AddOperator
	SubOperator
	MulOperator
	DivOperator
	NegateOperator
	unaryOperatorEnd
	ModOperator
	EqualsOperator
	NotEqualsOperator
	ShiftLeftOperator
	ShiftRightOperator
	AndOperator
	XorOperator
	OrOperator
	BitOrOperator
	BitAndOperator
	GreaterOperator
	GreaterEqualsOperator
	assignOperatorBegin
	AssignOperator
	AddAssignOperator
	SubAssignOperator
	MulAssignOperator
	DivAssignOperator
	assignOperatorEnd
	ColonOperator
	SmallerOperator
	SmallerEqualsOperator
	IncrementOperator
	DecrementOperator
	LeftParenOperator
	RightParenOperator
	LeftCurlyOperator
	RightCurlyOperator
	LeftBracketOperator
	RightBracketOperator
	SemicolonOperator
	CommaOperator
)

const InvalidOperatorName = "invalid"

var operatorNames = map[Operator]string{
	InvalidOperator:       InvalidOperatorName,
	AddOperator:           "+",
	SubOperator:           "-",
	MulOperator:           "*",
	DivOperator:           "/",
	ModOperator:           "%",
	EqualsOperator:        "==",
	NotEqualsOperator:     "!=",
	ShiftLeftOperator:     "<<",
	ShiftRightOperator:    ">>",
	AndOperator:           "&&",
	XorOperator:           "^",
	OrOperator:            "||",
	BitOrOperator:         "|",
	BitAndOperator:        "&",
	GreaterOperator:       ">",
	GreaterEqualsOperator: ">=",
	NegateOperator:        "!",
	AssignOperator:        "==",
	AddAssignOperator:     "+=",
	SubAssignOperator:     "-=",
	MulAssignOperator:     "*=",
	DivAssignOperator:     "/=",
	ColonOperator:         ":",
	SmallerOperator:       "<",
	SmallerEqualsOperator: "<=",
	IncrementOperator:     "++",
	DecrementOperator:     "--",
	LeftParenOperator:     "(",
	RightParenOperator:    ")",
	LeftCurlyOperator:     "{",
	RightCurlyOperator:    "}",
	LeftBracketOperator:   "]",
	RightBracketOperator:  "[",
	SemicolonOperator:     ";",
	CommaOperator:         ",",
}

type Precedence int8

const (
	LowPrecedence   = 0
	UnaryPrecedence = 6
	HighPrecedence  = 7
)

func (operator Operator) Precedence() Precedence {
	switch operator {
	case EqualsOperator,
		NotEqualsOperator,
		GreaterOperator,
		GreaterEqualsOperator:
		return 3
	case AddOperator,
		SubOperator,
		ModOperator,
		OrOperator,
		XorOperator:
		return 4
	case MulOperator,
		DivOperator,
		ShiftLeftOperator,
		ShiftRightOperator,
		AndOperator:
		return LowPrecedence
	}
	return 0
}

func (operator Operator) IsAssign() bool {
	return operator > assignOperatorBegin && operator < assignOperatorEnd
}

func (operator Operator) IsUnaryOperator() bool {
	return operator > unaryOperatorBegin && operator < unaryOperatorEnd
}

func (operator Operator) String() string {
	name, ok := operatorNames[operator]
	if !ok {
		return InvalidOperatorName
	}
	return name
}
