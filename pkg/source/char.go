package source

type Char rune

const (
	EndOfFile = Char(-1)
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

func (char Char) IsWhitespace() bool {
	return char == ' ' || char == '\t'
}
