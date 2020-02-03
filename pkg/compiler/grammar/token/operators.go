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
	QuestionMarkOperator
	CommaOperator
	DotOperator
)

const InvalidOperatorName = "invalid"

var operatorNames = map[Operator]string{
	InvalidOperator:       InvalidOperatorName,
	AddOperator:           "+",
	QuestionMarkOperator:  "?",
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

// Lists all supported precedences, ordered from weak to strong. Using this
// precedences, the parser can choose not to include certain operators, such
// as conditional ones, while still parsing arithmetic expressions. Tokens
// should only have precedence values that are listed here.
const (
	// LowestPrecedence is the precedence value for operators that may not
	// appear in a unary-expression. When the parser encounters an operator
	// which this precedence, it knows that the expression's end is reached.
	LowestPrecedence Precedence = iota
	InitialPrecedence
	// InitialConditionalPrecedence is the precedence that the parsers
	// starts at when parsing conditional expressions. It can also be used
	// when parsing arithmetic expressions. This value has to be greater than
	// LowestPrecedence and the conditional precedences.
	InitialConditionalPrecedence
	WeakLogicalPrecedence
	StrongLogicalPrecedence
	RelationalPrecedence
	// InitialArithmeticPrecedence is the precedence that the parser starts
	// at when parsing arithmetic expressions. It prevents inclusion of
	// conditional operators. This value has to be greater than
	// LowestPrecedence and the arithmetic precedences.
	InitialArithmeticPrecedence
	WeakArithmeticPrecedence
	StrongArithmeticPrecedence
)

func (precedence Precedence) Next() Precedence {
	return precedence + 1
}

var precedenceTable = map[Operator]Precedence{
	MulOperator:        StrongArithmeticPrecedence,
	DivOperator:        StrongArithmeticPrecedence,
	ModOperator:        StrongArithmeticPrecedence,
	ShiftLeftOperator:  StrongArithmeticPrecedence,
	ShiftRightOperator: StrongArithmeticPrecedence,

	AddOperator: WeakArithmeticPrecedence,
	SubOperator: WeakArithmeticPrecedence,
	XorOperator: WeakArithmeticPrecedence,

	SmallerEqualsOperator: RelationalPrecedence,
	SmallerOperator:       RelationalPrecedence,
	GreaterEqualsOperator: RelationalPrecedence,
	GreaterOperator:       RelationalPrecedence,
	EqualsOperator:        RelationalPrecedence,
	NotEqualsOperator:     RelationalPrecedence,

	AndOperator: StrongLogicalPrecedence,
	OrOperator:  WeakLogicalPrecedence,
}

func (operator Operator) Precedence() Precedence {
	if precedence, ok := precedenceTable[operator]; ok {
		return precedence
	}
	return LowestPrecedence
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
