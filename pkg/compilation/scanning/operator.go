package scanning

import (
	"errors"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
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

func (scanning *Scanning) findOperatorOption(options OperatorOptions, char source.Char) (token.Operator, error) {
	operator, ok := options[scanning.peekChar()]
	if ok {
		scanning.advance()
		return operator, nil
	}
	singleOperator, ok := options[singleChar]
	if !ok {
		return token.InvalidOperator, ErrNoSuchOperator
	}
	scanning.advance()
	return singleOperator, nil
}

func operatorOptionsOfChar(char source.Char) (OperatorOptions, bool) {
	options, ok := operatorTable[char]
	return options, ok
}
