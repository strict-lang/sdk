package scanner

import (
	"errors"
	"github.com/BenjaminNitschke/Strict/pkg/source"
	"github.com/BenjaminNitschke/Strict/pkg/token"
)

var (
	ErrNoSuchOperator = errors.New("there is no such operator")
)

type OperatorOptions map[source.Char]token.Operator
type OperatorTable map[source.Char]OperatorOptions

const singleChar = source.Char(0)

func singleOperatorOption(operator token.Operator) OperatorOptions {
	return OperatorOptions{singleChar: operator}
}

func isKnownOperator(char source.Char) bool {
	_, ok := operatorTable[char]
	return ok
}

func (scanner *Scanner) ScanOperator() token.Token {
	operator, err := scanner.gatherOperator()
	if err != nil {
		scanner.reportError(err)
		return scanner.createInvalidToken()
	}
	return token.NewOperatorToken(operator, scanner.currentPosition(), scanner.indent)
}

func (scanner Scanner) gatherOperator() (token.Operator, error) {
	peek := scanner.reader.Pull()
	options, ok := operatorOptionsOfChar(peek)
	if !ok || len(options) == 0 {
		return token.InvalidOperator, ErrNoSuchOperator
	}
	next := scanner.reader.Peek()
	return scanner.findOperatorOption(options, next)
}

func (scanner *Scanner) findOperatorOption(
	options OperatorOptions,
	char source.Char) (token.Operator, error) {

	operator, ok := options[char]
	if ok {
		scanner.reader.Pull()
		return operator, nil
	}
	singleOperator, ok := options[singleChar]
	if !ok {
		return token.InvalidOperator, ErrNoSuchOperator
	}
	return singleOperator, nil
}

func operatorOptionsOfChar(char source.Char) (OperatorOptions, bool) {
	options, ok := operatorTable[char]
	return options, ok
}
