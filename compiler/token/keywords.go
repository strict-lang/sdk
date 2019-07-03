package token

type Keyword uint8

const (
	KeywordTokenName   = "keyword"
	InvalidKeywordName = "invalid"
)

const (
	InvalidKeyword Keyword = iota
	MethodKeyword
	TypeKeyword
	IsKeyword
	IfKeyword
	ElseKeyword
	YieldKeyword
	ReturnKeyword
	ToKeyword
	OrKeyword
	NotKeyword
	AndKeyword
	ForKeyword
)

var keywordNameTable = map[Keyword]string{
	MethodKeyword: "method",
	ReturnKeyword: "return",
	TypeKeyword:   "type",
	IsKeyword:     "is",
	IfKeyword:     "if",
	ElseKeyword:   "else",
	ForKeyword:    "for",
	YieldKeyword:  "yield",
	ToKeyword:     "to",
	AndKeyword:    "and",
	OrKeyword:     "or",
	NotKeyword:    "not",
}

var operatorKeywords = map[Keyword]Operator{
	OrKeyword:  OrOperator,
	AndKeyword: AndOperator,
}

var keywordNameLookupTable map[string]Keyword

func init() {
	keywordNameLookupTable = make(map[string]Keyword)
	for keyword, name := range keywordNameTable {
		keywordNameLookupTable[name] = keyword
	}
}

func KeywordByName(name string) (Keyword, bool) {
	keyword, ok := keywordNameLookupTable[name]
	return keyword, ok
}

func KeywordNames() map[string]Keyword {
	return keywordNameLookupTable
}

func (keyword Keyword) String() string {
	name, ok := keywordNameTable[keyword]
	if !ok {
		return InvalidKeywordName
	}
	return name
}