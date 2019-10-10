package lexical

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/input"
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

func findEscapedCharacter(char input.Char) (input.Char, bool) {
	if escaped, ok := escapedCharacters[rune(char)]; ok {
		return input.Char(escaped), true
	}
	return input.EndOfFile, false
}
