package source

type Char rune

// Reserved special characters that are not part of any charset and only used internally to
// indicate the absence or position of a character in a stream.
const (
	EndOfFile   = Char(-1)
	BeginOfFile = Char(-1)
	NoChar = Char(-2)
)

func (char Char) IsInRange(lower, upper Char) bool {
	return lower <= char && char <= upper
}

func (char Char) IsAlphabetic() bool {
	return char.IsInRange('a', 'z') || char.IsInRange('A', 'Z')
}

func (char Char) IsNumeric() bool {
	return char.IsInRange('0', '9')
}

func (char Char) IsAlphanumeric() bool {
	return char.IsNumeric() || char.IsAlphabetic()
}

func (char Char) IsWhitespaceOrLineFeed() bool {
	return char.IsWhitespace() || char.IsLineFeed()
}

func (char Char) IsWhitespace() bool {
	return char == ' ' || char == '\t'
}

func (char Char) IsLineFeed() bool {
	return char == '\n' || char == '\r'
}

func (char Char) DigitValue() int {
	switch {
	case char.IsInRange('0', '9'):
		return int(char - '0')
	case char.IsInRange('a', 'f'):
		return int(char - 'a' + 10)
	case char.IsInRange('A', 'F'):
		return int(char - 'A' + 10)
	}
	return 16
}
