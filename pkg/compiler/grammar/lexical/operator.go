package lexical

import (
	"errors"
	"gitlab.com/strict-lang/sdk/pkg/compilation/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compilation/input"
)

var (
	ErrNoSuchOperator = errors.New("there is no such operator")
)

type OperatorOptions map[input.Char]token.Operator
type OperatorTable map[input.Char]OperatorOptions

const singleChar = input.Char(0)

func singleOperatorOption(operator token.Operator) OperatorOptions {
	return OperatorOptions{singleChar: operator}
}

func isKnownOperator(char input.Char) bool {
	_, ok := operatorTable[char]
	return ok
}

func (scanning *Scanning) ScanOperator() token.Token {
	operator, err := scanning.gatherOperator()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	scanning.operatorGathered(operator)
	return token.NewOperatorToken(operator, scanning.currentPosition(), scanning.indent)
}

// operatorGathered is called once an operator is gathered. It checks whether the operator
// is enabling or disabling the 'insertEos' flag and applies it.
func (scanning *Scanning) operatorGathered(operator token.Operator) {
	if _, ok := endOfStatementDisablingOperators[operator]; ok {
		scanning.endOfStatementPrevention++
		return
	}
	if _, ok := endOfStatementEnablingOperators[operator]; ok {
		scanning.endOfStatementPrevention--
	}
}

func (scanning Scanning) gatherOperator() (token.Operator, error) {
	char := scanning.char()
	options, ok := operatorOptionsOfChar(char)
	if !ok || len(options) == 0 {
		return token.InvalidOperator, ErrNoSuchOperator
	}
	return scanning.findOperatorOption(options, scanning.char())
}

func (scanning *Scanning) findOperatorOption(options OperatorOptions, char input.Char) (token.Operator, error) {
	operator, ok := options[scanning.peekChar()]
	if ok {
		scanning.advance()
		scanning.advance() // Advance twice since first char is still current
		return operator, nil
	}
	singleOperator, ok := options[singleChar]
	if !ok {
		return token.InvalidOperator, ErrNoSuchOperator
	}
	scanning.advance()
	return singleOperator, nil
}

func operatorOptionsOfChar(char input.Char) (OperatorOptions, bool) {
	options, ok := operatorTable[char]
	return options, ok
}
