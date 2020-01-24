package token

import "fmt"

type KeywordToken struct {
	Keyword  Keyword
	position Position
	indent   Indent
}

func NewKeywordToken(keyword Keyword, position Position, indent Indent) Token {
	return &KeywordToken{
		Keyword:  keyword,
		position: position,
		indent:   indent,
	}
}

func (keyword KeywordToken) String() string {
	return fmt.Sprintf("%s(%s)", KeywordTokenName, keyword.Keyword)
}

func (keyword KeywordToken) Name() string {
	return KeywordTokenName
}

func (keyword KeywordToken) Value() string {
	return keyword.Keyword.String()
}

func (keyword KeywordToken) Position() Position {
	return keyword.position
}

func (keyword KeywordToken) Indent() Indent {
	return keyword.indent
}

func (keyword KeywordToken) IsOperatorKeyword() bool {
	_, ok := operatorKeywords[keyword.Keyword]
	return ok
}

func (keyword KeywordToken) AsOperator() Operator {
	operator, ok := operatorKeywords[keyword.Keyword]
	if !ok {
		return InvalidOperator
	}
	return operator
}

func IsKeywordToken(token Token) bool {
	_, ok := token.(*KeywordToken)
	return ok
}

func HasOperatorValue(token Token, value Operator) bool {
	return OperatorValue(token) == value
}

func HasKeywordValue(token Token, value Keyword) bool {
	keyword, ok := token.(*KeywordToken)
	if !ok {
		return false
	}
	return keyword.Keyword == value
}

func IsAssignOperator(token Token) bool {
	return OperatorValue(token).IsAssign()
}

func IsOperatorOrOperatorKeywordToken(token Token) bool {
	if _, ok := token.(*OperatorToken); ok {
		return true
	}
	keyword, ok := token.(*KeywordToken)
	if !ok {
		return false
	}
	return keyword.IsOperatorKeyword()
}

func OperatorValue(token Token) Operator {
	if operator, ok := token.(*OperatorToken); ok {
		return operator.Operator
	}
	keyword, ok := token.(*KeywordToken)
	if !ok {
		return InvalidOperator
	}
	return keyword.AsOperator()
}

func KeywordValue(token Token) Keyword {
	keyword, ok := token.(*KeywordToken)
	if !ok {
		return InvalidKeyword
	}
	return keyword.Keyword
}

func KeywordValueOfOperator(operator Operator) (Keyword, bool) {
	keyword, ok := operatorKeywordsReversed[operator]
	return keyword, ok
}
