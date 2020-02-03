package lexical

import (
	"errors"
	"strict.dev/sdk/pkg/compiler/grammar/token"
	"strict.dev/sdk/pkg/compiler/input"
)

var (
	errNoSuchOperator = errors.New("there is no such operator")
)

type operatorOptions map[input.Char]token.Operator
type operatorOptionTable map[input.Char]operatorOptions

const singleChar = input.Char(0)

func singleOperatorOption(operator token.Operator) operatorOptions {
	return operatorOptions{singleChar: operator}
}

func isKnownOperator(char input.Char) bool {
	_, ok := operatorTable[char]
	return ok
}

func (scanning *Scanning) scanOperator() token.Token {
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
	options, ok := listOperatorOptions(char)
	if !ok || len(options) == 0 {
		return token.InvalidOperator, errNoSuchOperator
	}
	return scanning.findOperatorOption(options, scanning.char())
}

func (scanning *Scanning) findOperatorOption(
	options operatorOptions, char input.Char) (token.Operator, error) {

	operator, ok := options[scanning.peekChar()]
	if ok {
		scanning.advance()
		scanning.advance() // Advance twice since first char is still current
		return operator, nil
	}
	singleOperator, ok := options[singleChar]
	if !ok {
		return token.InvalidOperator, errNoSuchOperator
	}
	scanning.advance()
	return singleOperator, nil
}

func listOperatorOptions(char input.Char) (operatorOptions, bool) {
	options, ok := operatorTable[char]
	return options, ok
}
