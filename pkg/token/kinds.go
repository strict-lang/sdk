package token

type Kind uint8

const (
	Invalid Kind = iota
	Identifier
	StringLiteral
	keywordsBegin
	MethodKeyword
	keywordsEnd
	operatorsBegin
	Assign
	SmallerThan
	GreaterThan
	SmallerEqual
	GreaterEqual
	Or
	And
	Is
	IsNot
	Divide
	Multiply
	Modulo
	Not
	operatorsEnd
)

const invalidName = "invalid"

var names = map[Kind]string{
	Invalid:       invalidName,
	Identifier:    "identifier",
	MethodKeyword: "method",
	Assign:        "=",
	SmallerThan:   "<",
	GreaterThan:   ">",
	SmallerEqual:  "<=",
	GreaterEqual:  ">=",
	Or:            "or",
	Is:            "is",
	IsNot:         "is not",
	Divide:        "/",
	Multiply:      "*",
	Modulo:        "%",
	Not:           "not",
}

func (kind Kind) isInExclusiveRange(lower, upper Kind) bool {
	return lower <= kind && kind <= upper
}

func (kind Kind) IsKeyword() bool {
	return kind.isInExclusiveRange(keywordsBegin, keywordsEnd)
}

func (kind Kind) IsOperator() bool {
	return kind.isInExclusiveRange(operatorsBegin, operatorsEnd)
}

func (kind Kind) Group() string {
	if kind.IsOperator() {
		return "Operator"
	}
	if kind.IsKeyword() {
		return "Keyword"
	}
	return NameOfKind(kind)
}

func (kind Kind) Name() string {
	return NameOfKind(kind)
}

func NameOfKind(kind Kind) string {
	name, ok := names[kind]
	if !ok {
		return invalidName
	}
	return name
}
