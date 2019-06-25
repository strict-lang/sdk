package token

type Keyword uint8

const (
	KeywordTokenName = "keyword"
	InvalidKeywordName = "invalid"
)

const (
	MethodKeyword Keyword = iota
	TypeKeyword
)

var keywordNameTable = map[Keyword] string{
	MethodKeyword: "method",
	TypeKeyword: "type",
}

var keywordNameLookupTable map[string] Keyword

func init() {
	keywordNameLookupTable = make(map[string] Keyword)
	for keyword, name := range keywordNameTable {
		keywordNameLookupTable[name] = keyword
	}
}

type KeywordToken struct {
	keyword Keyword
	position Position
}

func NewKeywordToken(keyword Keyword, position Position) Token {
	return &KeywordToken{
		keyword: keyword,
		position: position,
	}
}

func (keyword KeywordToken) String() string {
	name, ok := keywordNameTable[keyword.keyword]
	if !ok {
		return InvalidKeywordName
	}
	return name
}

func (keyword KeywordToken) Name() string {
	return KeywordTokenName
}

func (keyword KeywordToken) Value() string {
	return string(keyword.keyword)
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
