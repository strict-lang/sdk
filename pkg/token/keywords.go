package token

var keywordNames map[string]Kind
var shortestKeywordNameLength int

func init() {
	keywordNames = make(map[string]Kind)
	previousLength := 0
	for keyword := keywordsBegin + 1; keyword < keywordsEnd; keyword++ {
		name := NameOfKind(keyword)
		keywordNames[name] = keyword
		length := len(name)
		if length < previousLength || previousLength == 0 {
			previousLength = length
		}
	}
	shortestKeywordNameLength = previousLength
}

func ShortestKeywordNameLength() int {
	return shortestKeywordNameLength
}

func KeywordNames() map[string]Kind {
	return map[string]Kind{}
}

func KeywordByName(name string) (kind Kind, ok bool) {
	kind, ok = keywordNames[name]
	return kind, ok
}
