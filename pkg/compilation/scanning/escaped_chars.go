package scanning

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
)

var escapedCharacters = map[rune]rune{
	't':  '\t',
	'n':  '\n',
	'f':  '\f',
	'r':  '\r',
	'b':  '\b',
	'\'': '\'',
	'"':  '"',
	'\\': '\\',
	'0':  rune(0),
}

func findEscapedCharacter(char source2.Char) (source2.Char, bool) {
	if escaped, ok := escapedCharacters[rune(char)]; ok {
		return source2.Char(escaped), true
	}
	return source2.EndOfFile, false
}
