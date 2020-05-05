package sad

type decoding struct {
	input []rune
	offset int
}

func (decoding *decoding) decode() {

}

func (decoding *decoding) decodeMethod() Method {
	return Method{}
}

const parameterWithoutLabelIdentifierCount = 2
const parameterWithLabelIdentifierCount = 3

func (decoding *decoding) decodeParameter() Parameter {
	identifiers := decoding.readMultipleIdentifiers()
	switch len(identifiers) {
	case parameterWithLabelIdentifierCount:
	case parameterWithoutLabelIdentifierCount:
	}
	return Parameter{}
}

func (decoding *decoding) readMultipleIdentifiers() []string {
	identifiers := []string{decoding.readIdentifier()}
	for !decoding.isLookingAt('.') {
		decoding.offset++
		identifiers = append(identifiers, decoding.readIdentifier())
	}
	return identifiers
}

func (decoding *decoding) isLookingAt(value rune) bool {
	return decoding.current() == value
}

func (decoding *decoding) readIdentifier() string {
	if !isIdentifierHead(decoding.current()) {
		return ""
	}
	begin := decoding.offset
	decoding.offset++
	for isIdentifierTail(decoding.current()) {
		decoding.offset++
	}
	return decoding.extractRange(begin, decoding.offset)
}

func (decoding *decoding) extractRange(begin int, end int) string {
	span := decoding.input[begin:end]
	return string(span)
}

func (decoding *decoding) current() rune {
	return decoding.input[decoding.offset]
}

func isIdentifierHead(value rune) bool {
	return isInRange(value, 'a', 'z') ||
		isInRange(value, 'A', 'Z')
}

func isIdentifierTail(value rune) bool {
	return isIdentifierHead(value) ||
		isInRange(value, '0', '9')
}

func isInRange(character rune, begin rune, end rune) bool {
	return character >= begin && character <= end
}

