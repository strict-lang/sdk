package scanning

import (
	"errors"
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

var (
	ErrNoSuchOperator = errors.New("there is no such operator")
)

type OperatorOptions map[source2.Char]token2.Operator
type OperatorTable map[source2.Char]OperatorOptions

const singleChar = source2.Char(0)

func singleOperatorOption(operator token2.Operator) OperatorOptions {
	return OperatorOptions{singleChar: operator}
}

func isKnownOperator(char source2.Char) bool {
	_, ok := operatorTable[char]
	return ok
}

func (scanning *Scanning) ScanOperator() token2.Token {
	operator, err := scanning.gatherOperator()
	if err != nil {
		scanning.reportError(err)
		return scanning.createInvalidToken()
	}
	scanning.operatorGathered(operator)
	return token2.NewOperatorToken(operator, scanning.currentPosition(), scanning.indent)
}

// operatorGathered is called once an operator is gathered. It checks whether the operator
// is enabling or disabling the 'insertEos' flag and applies it.
func (scanning *Scanning) operatorGathered(operator token2.Operator) {
	if _, ok := endOfStatementDisablingOperators[operator]; ok {
		scanning.endOfStatementPrevention++
		return
	}
	if _, ok := endOfStatementEnablingOperators[operator]; ok {
		scanning.endOfStatementPrevention--
	}
}

func (scanning Scanning) gatherOperator() (token2.Operator, error) {
	char := scanning.char()
	options, ok := operatorOptionsOfChar(char)
	if !ok || len(options) == 0 {
		return token2.InvalidOperator, ErrNoSuchOperator
	}
	return scanning.findOperatorOption(options, scanning.char())
}

func (scanning *Scanning) findOperatorOption(options OperatorOptions, char source2.Char) (token2.Operator, error) {
	operator, ok := options[scanning.peekChar()]
	if ok {
		scanning.advance()
		return operator, nil
	}
	singleOperator, ok := options[singleChar]
	if !ok {
		return token2.InvalidOperator, ErrNoSuchOperator
	}
	scanning.advance()
	return singleOperator, nil
}

func operatorOptionsOfChar(char source2.Char) (OperatorOptions, bool) {
	options, ok := operatorTable[char]
	return options, ok
}
