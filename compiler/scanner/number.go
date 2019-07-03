package scanner

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/token"
	"strings"
)

type Radix int8

const (
	Binary      Radix = 2
	Decimal     Radix = 10
	Hexadecimal Radix = 16
)

func (scanner *Scanner) ScanNumber() token.Token {
	number, err := scanner.gatherNumber()
	if err != nil {
		scanner.reportError(err)
		return scanner.createInvalidToken()
	}
	return token.NewNumberLiteralToken(number, scanner.currentPosition(), scanner.indent)
}

func (scanner *Scanner) gatherNumericDigits(builder *strings.Builder, base Radix) {
	for scanner.reader.Last().DigitValue() < int(base) {
		if scanner.reader.IsExhausted() {
			return
		}
		builder.WriteRune(rune(scanner.reader.Last()))
		scanner.reader.Pull()
	}
}

func (scanner *Scanner) gatherNumber() (string, error) {
	var builder strings.Builder
	switch scanner.reader.Pull() {
	case '0':
		builder.WriteRune('0')
		switch scanner.reader.Pull() {
		case 'x', 'X':
			builder.WriteRune('x')
			scanner.reader.Pull()
			return scanner.gatherNumberWithRadix(&builder, Hexadecimal)
		case 'b', 'B':
			builder.WriteRune('b')
			scanner.reader.Pull()
			return scanner.gatherNumberWithRadix(&builder, Binary)
		case '.':
			builder.WriteRune('.')
			scanner.reader.Pull()
			err := scanner.gatherFloatingPointNumber(&builder)
			return builder.String(), err
		}
	}
	scanner.gatherNumericDigits(&builder, Decimal)
	if scanner.reader.Last() == '.' {
		builder.WriteRune('.')
		scanner.reader.Pull()
		if err := scanner.gatherFloatingPointNumber(&builder); err != nil {
			return scanner.reader.String(), err
		}
	}
	return builder.String(), nil
}

func (scanner *Scanner) gatherExponent(builder *strings.Builder) error {
	switch scanner.reader.Last() {
		case '-', '+':
			if scanner.reader.Pull().DigitValue() < int(Decimal) {
				scanner.gatherNumericDigits(builder, Decimal)
				return nil
			}
			return &UnexpectedCharError{
				got: scanner.reader.Last(),
				expected: "digital number",
			}
	}
	return nil
}

func (scanner *Scanner) gatherFloatingPointNumber(builder *strings.Builder) error {
	scanner.gatherNumericDigits(builder, Decimal)
	switch scanner.reader.Last() {
	case 'e', 'E':
		builder.WriteRune('e')
		scanner.reader.Pull()
		return scanner.gatherExponent(builder)
	}
	return nil
}

func (scanner *Scanner) gatherNumberWithRadix(builder *strings.Builder, radix Radix) (string, error) {
	scanner.gatherNumericDigits(builder, radix)
	if scanner.reader.Last().DigitValue() >= int(radix) {
		return scanner.reader.String(), &UnexpectedCharError{
			got: scanner.reader.Last(),
			expected: fmt.Sprintf("number with radix %d", radix),
		}
	}
	return builder.String(), nil
}