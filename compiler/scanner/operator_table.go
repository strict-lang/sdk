package scanner

import "github.com/BenjaminNitschke/Strict/compiler/token"

var operatorTable = OperatorTable{
	'+': {
		singleChar: token.AddOperator,
		'=':        token.AddAssignOperator,
		'+':        token.IncrementOperator,
	},
	'-': {
		singleChar: token.SubOperator,
		'=':        token.SubAssignOperator,
		'-':        token.DecrementOperator,
	},
	'/': {
		singleChar: token.DivOperator,
		'=':        token.DivAssignOperator,
	},
	'*': {
		singleChar: token.MulOperator,
		'=':        token.MulAssignOperator,
	},
	'=': {
		singleChar: token.AssignOperator,
		'=':        token.EqualsOperator,
	},
	'!': {
		singleChar: token.NegateOperator,
		'=':        token.NotEqualsOperator,
	},
	'|': {
		singleChar: token.BitOrOperator,
		'|':        token.OrOperator,
	},
	'&': {
		singleChar: token.BitAndOperator,
		'&':        token.AndOperator,
	},
	'>': {
		singleChar: token.GreaterOperator,
		'=':        token.GreaterEqualsOperator,
		// '>':        token.ShiftRightOperator,
	},
	'<': {
		singleChar: token.SmallerOperator,
		'=':        token.SmallerEqualsOperator,
		// '<':        token.ShiftLeftOperator,
	},
	';': singleOperatorOption(token.SemicolonOperator),
	':': singleOperatorOption(token.ColonOperator),
	')': singleOperatorOption(token.RightParenOperator),
	'(': singleOperatorOption(token.LeftParenOperator),
	'{': singleOperatorOption(token.LeftCurlyOperator),
	'}': singleOperatorOption(token.RightCurlyOperator),
	'[': singleOperatorOption(token.LeftBracketOperator),
	']': singleOperatorOption(token.RightBracketOperator),
	',': singleOperatorOption(token.CommaOperator),
}

// endOfStatementDisablingOperators are operators that disable the scanners 'insertEos' flag.
// If the scanner gathers one of those operators, it changes the flag to false. The maps
// keys are the disabling operators and their values are the corresponding enabling operators.
var endOfStatementDisablingOperators = map[token.Operator]token.Operator{
	token.LeftParenOperator:   token.RightParenOperator,
	token.LeftBracketOperator: token.RightBracketOperator,
}

// endOfStatementEnablingOperators is a reversed map of the endOfStatementDisablingOperators.
var endOfStatementEnablingOperators map[token.Operator]token.Operator

func init() {
	length := len(endOfStatementDisablingOperators)
	endOfStatementEnablingOperators = make(map[token.Operator]token.Operator, length)
	for disabler, enabler := range endOfStatementDisablingOperators {
		endOfStatementEnablingOperators[enabler] = disabler
	}
}
