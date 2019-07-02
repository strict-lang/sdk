package scanner

import "github.com/BenjaminNitschke/Strict/pkg/token"

type Radix int8

const (
	DecimalRadix     Radix = iota
	HexadecimalRadix Radix = iota
	BinaryRadix      Radix = iota
)

func (scanner *Scanner) ScanNumber() token.Token {
	return token.NewInvalidToken("", token.Position{}, scanner.indent)
}

func (scanner *Scanner) GatherDecimal() (string, error) {
	return "", nil
}

func (scanner *Scanner) GatherHexadecimal() (string, error) {
	return "", nil
}

func (scanner *Scanner) GatherBinary() (string, error) {
	return "", nil
}

func (scanner *Scanner) GatherFloat() (string, error) {
	return "", nil
}
