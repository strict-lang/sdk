package scanner

import "github.com/BenjaminNitschke/Strict/compiler/source"

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

func findEscapedCharacter(char source.Char) (source.Char, bool) {
	if escaped, ok := escapedCharacters[rune(char)]; ok {
		return source.Char(escaped), true
	}
	return source.EndOfFile, false
}
