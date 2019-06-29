package token

type Keyword uint8

const (
	KeywordTokenName   = "keyword"
	InvalidKeywordName = "invalid"
)

const (
	MethodKeyword Keyword = iota
	TypeKeyword
)

var keywordNameTable = map[Keyword]string{
	MethodKeyword: "method",
	TypeKeyword:   "type",
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
