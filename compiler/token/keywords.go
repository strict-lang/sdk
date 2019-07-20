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
	InKeyword
	IfKeyword
	ElseKeyword
	YieldKeyword
	ReturnKeyword
	FromKeyword
	ToKeyword
	DoKeyword
	OrKeyword
	IsNotKeyword
	AndKeyword
	ForKeyword
)

var keywordNameTable = map[Keyword]string{
	MethodKeyword: "method",
	ReturnKeyword: "return",
	TypeKeyword:   "type",
	IsKeyword:     "is",
	InKeyword:     "in",
	IfKeyword:     "if",
	ElseKeyword:   "else",
	ForKeyword:    "for",
	FromKeyword:   "from",
	YieldKeyword:  "yield",
	ToKeyword:     "to",
	DoKeyword:     "do",
	AndKeyword:    "and",
	OrKeyword:     "or",
	IsNotKeyword:  "isnt",
}

var operatorKeywords = map[Keyword]Operator{
	OrKeyword:    OrOperator,
	AndKeyword:   AndOperator,
	IsKeyword:    EqualsOperator,
	IsNotKeyword: NotEqualsOperator,
}

var operatorKeywordsReversed map[Operator]Keyword

var keywordNameLookupTable map[string]Keyword

func init() {
	keywordNameLookupTable = make(map[string]Keyword)
	for keyword, name := range keywordNameTable {
		keywordNameLookupTable[name] = keyword
	}
	operatorKeywordsReversed = make(map[Operator]Keyword)
	for keyword, operator := range operatorKeywords {
		operatorKeywordsReversed[operator] = keyword
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
