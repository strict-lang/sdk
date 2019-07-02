package token

type KeywordToken struct {
	Keyword  Keyword
	position Position
	indent Indent
}

func NewKeywordToken(keyword Keyword, position Position, indent Indent) Token {
	return &KeywordToken{
		Keyword:  keyword,
		position: position,
		indent: indent,
	}
}

func (keyword KeywordToken) String() string {
	return keyword.Keyword.String()
}

func (keyword KeywordToken) Name() string {
	return KeywordTokenName
}

func (keyword KeywordToken) Value() string {
	return string(keyword.Keyword)
}

func (keyword KeywordToken) Position() Position {
	return keyword.position
}

func (KeywordToken) IsKeyword() bool {
	return true
}

func (KeywordToken) IsOperator() bool {
	return false
}

func (KeywordToken) IsLiteral() bool {
	return false
}

func (KeywordToken) IsValid() bool {
	return false
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
