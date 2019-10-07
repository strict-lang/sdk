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
	ArrowOperator
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
	DotOperator
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
	AssignOperator:        "=",
	AddAssignOperator:     "+=",
	SubAssignOperator:     "-=",
	MulAssignOperator:     "*=",
	DivAssignOperator:     "/=",
	ColonOperator:         ":",
	ArrowOperator:         "=>",
	SmallerOperator:       "<",
	SmallerEqualsOperator: "<=",
	IncrementOperator:     "++",
	DecrementOperator:     "--",
	LeftParenOperator:     "(",
	RightParenOperator:    ")",
	LeftCurlyOperator:     "{",
	RightCurlyOperator:    "}",
	LeftBracketOperator:   "[",
	RightBracketOperator:  "]",
	SemicolonOperator:     ";",
	CommaOperator:         ",",
	DotOperator:           ".",
}

type Precedence int8

const (
	LowPrecedence   = 0
	UnaryPrecedence = 7
	HighPrecedence  = 8
)

var precedenceTable = map[Operator] Precedence {
	MulOperator: 5,
	DivOperator: 5,
	ModOperator: 5,
	ShiftLeftOperator: 5,
	ShiftRightOperator: 5,

	AddOperator: 4,
	SubOperator: 4,
	XorOperator: 4,
	SmallerEqualsOperator: 3,
	SmallerOperator: 3,
	GreaterEqualsOperator: 3,
	GreaterOperator: 3,
	EqualsOperator: 3,
	NotEqualsOperator: 3,
	AndOperator: 2,
	OrOperator: 1,
}

func (operator Operator) Precedence() Precedence {
	if precedence, ok := precedenceTable[operator]; ok {
		return precedence
	}
	return LowPrecedence
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
