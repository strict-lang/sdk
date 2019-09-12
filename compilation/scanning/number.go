package scanning

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/token"
	"strings"
)

type Radix int8

const (
	Binary      Radix = 2
	Decimal     Radix = 10
	Hexadecimal Radix = 16
)

func (scanning *Scanning) ScanNumber() token.Token {
	number, err := scanning.gatherNumber()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token.NewNumberLiteralToken(number, scanning.currentPosition(), scanning.indent)
}

func isDigitInRadix(digitValue int, base Radix) bool {
	return digitValue < int(base)
}

func (scanning *Scanning) gatherNumericDigits(builder *strings.Builder, base Radix) {
	digitValue := scanning.reader.Last().DigitValue()
	if !isDigitInRadix(digitValue, base) {
		return
	}
	builder.WriteRune(rune(scanning.reader.Last()))
	for {
		char := scanning.reader.Pull()
		if !isDigitInRadix(char.DigitValue(), base) {
			return
		}
		builder.WriteRune(rune(char))
	}
}

func (scanning *Scanning) gatherNumber() (string, error) {
	scanning.reader.Pull()
	var builder strings.Builder
	if scanning.reader.Last() == '0' {
		builder.WriteRune('0')
		scanning.reader.Pull()
		switch scanning.reader.Last() {
		case 'x', 'X':
			builder.WriteRune('x')
			scanning.reader.Pull()
			return scanning.gatherNumberWithRadix(&builder, Hexadecimal)
		case 'b', 'B':
			builder.WriteRune('b')
			scanning.reader.Pull()
			return scanning.gatherNumberWithRadix(&builder, Binary)
		case '.':
			builder.WriteRune('.')
			scanning.reader.Pull()
			err := scanning.gatherFloatingPointNumber(&builder)
			return builder.String(), err
		default:
			return builder.String(), nil
		}
	}
	scanning.gatherNumericDigits(&builder, Decimal)
	if scanning.reader.Last() == '.' {
		builder.WriteRune('.')
		scanning.reader.Pull()
		if err := scanning.gatherFloatingPointNumber(&builder); err != nil {
			return scanning.reader.String(), err
		}
	}
	return builder.String(), nil
}

func (scanning *Scanning) gatherExponent(builder *strings.Builder) error {
	switch scanning.reader.Last() {
	case '-', '+':
		if scanning.reader.Pull().DigitValue() < int(Decimal) {
			scanning.gatherNumericDigits(builder, Decimal)
			return nil
		}
		return &UnexpectedCharError{
			got:      scanning.reader.Last(),
			expected: "digital number",
		}
	}
	return nil
}

func (scanning *Scanning) gatherFloatingPointNumber(builder *strings.Builder) error {
	scanning.gatherNumericDigits(builder, Decimal)
	switch scanning.reader.Last() {
	case 'e', 'E':
		builder.WriteRune('e')
		scanning.reader.Pull()
		return scanning.gatherExponent(builder)
	}
	return nil
}

func (scanning *Scanning) gatherNumberWithRadix(builder *strings.Builder, radix Radix) (string, error) {
	scanning.gatherNumericDigits(builder, radix)
	if scanning.reader.Last().DigitValue() >= int(radix) {
		return scanning.reader.String(), &UnexpectedCharError{
			got:      scanning.reader.Last(),
			expected: fmt.Sprintf("number with radix %d", radix),
		}
	}
	return builder.String(), nil
}
