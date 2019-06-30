package token

type KeywordToken struct {
	Keyword  Keyword
	position Position
}

func NewKeywordToken(keyword Keyword, position Position) Token {
	return &KeywordToken{
		Keyword:  keyword,
		position: position,
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
