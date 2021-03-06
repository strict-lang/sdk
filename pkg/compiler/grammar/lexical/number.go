package lexical

import (
	"fmt"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"strings"
)

type radix int8

const (
	binaryRadix      radix = 2
	decimalRadix     radix = 10
	hexadecimalRadix radix = 16
)

func (scanning *Scanning) scanNumber() token.Token {
	number, err := scanning.gatherNumber()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	return token.NewNumberLiteralToken(number, scanning.currentPosition(), scanning.indent)
}

func isDigitInRadix(digitValue int, base radix) bool {
	return digitValue < int(base)
}

func (scanning *Scanning) gatherNumericDigits(builder *strings.Builder, base radix) {
	digitValue := scanning.char().DigitValue()
	if !isDigitInRadix(digitValue, base) {
		return
	}
	builder.WriteRune(rune(scanning.char()))
	for {
		scanning.advance()
		char := scanning.char()
		if scanning.input.IsExhausted() {
			return
		}
		if !isDigitInRadix(char.DigitValue(), base) {
			return
		}
		builder.WriteRune(rune(char))
	}
}

func (scanning *Scanning) gatherNumber() (string, error) {
	var builder strings.Builder
	if scanning.char() == '0' {
		builder.WriteRune('0')
		scanning.advance()
		switch scanning.char() {
		case 'x', 'X':
			builder.WriteRune('x')
			scanning.advance()
			return scanning.gatherNumberWithRadix(&builder, hexadecimalRadix)
		case 'b', 'B':
			builder.WriteRune('b')
			scanning.advance()
			return scanning.gatherNumberWithRadix(&builder, binaryRadix)
		case '.':
			builder.WriteRune('.')
			scanning.advance()
			err := scanning.gatherFloatingPointNumber(&builder)
			return builder.String(), err
		default:
			return builder.String(), nil
		}
	}
	scanning.gatherNumericDigits(&builder, decimalRadix)
	if scanning.char() == '.' {
		builder.WriteRune('.')
		scanning.advance()
		if err := scanning.gatherFloatingPointNumber(&builder); err != nil {
			return scanning.input.String(), err
		}
	}
	return builder.String(), nil
}

func (scanning *Scanning) gatherExponent(builder *strings.Builder) error {
	switch scanning.char() {
	case '-', '+':
		scanning.advance()
		if scanning.char().DigitValue() < int(decimalRadix) {
			scanning.gatherNumericDigits(builder, decimalRadix)
			return nil
		}
		return &unexpectedCharError{
			got:      scanning.char(),
			expected: "digital number",
		}
	}
	return nil
}

func (scanning *Scanning) gatherFloatingPointNumber(builder *strings.Builder) error {
	scanning.gatherNumericDigits(builder, decimalRadix)
	switch scanning.char() {
	case 'e', 'E':
		builder.WriteRune('e')
		scanning.advance()
		return scanning.gatherExponent(builder)
	}
	return nil
}

func (scanning *Scanning) gatherNumberWithRadix(builder *strings.Builder, radix radix) (string, error) {
	scanning.gatherNumericDigits(builder, radix)
	if scanning.char().DigitValue() >= int(radix) {
		return scanning.input.String(), &unexpectedCharError{
			got:      scanning.char(),
			expected: fmt.Sprintf("number with radix %d", radix),
		}
	}
	scanning.advance()
	return builder.String(), nil
}
