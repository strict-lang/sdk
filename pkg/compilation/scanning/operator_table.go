package scanning

import (
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

var operatorTable = OperatorTable{
	'+': {
		singleChar: token2.AddOperator,
		'=':        token2.AddAssignOperator,
		'+':        token2.IncrementOperator,
	},
	'-': {
		singleChar: token2.SubOperator,
		'=':        token2.SubAssignOperator,
		'-':        token2.DecrementOperator,
	},
	'/': {
		singleChar: token2.DivOperator,
		'=':        token2.DivAssignOperator,
	},
	'*': {
		singleChar: token2.MulOperator,
		'=':        token2.MulAssignOperator,
	},
	'=': {
		singleChar: token2.AssignOperator,
		'=':        token2.EqualsOperator,
		'>':        token2.ArrowOperator,
	},
	'!': {
		singleChar: token2.NegateOperator,
		'=':        token2.NotEqualsOperator,
	},
	'|': {
		singleChar: token2.BitOrOperator,
		'|':        token2.OrOperator,
	},
	'&': {
		singleChar: token2.BitAndOperator,
		'&':        token2.AndOperator,
	},
	'>': {
		singleChar: token2.GreaterOperator,
		'=':        token2.GreaterEqualsOperator,
	},
	'<': {
		singleChar: token2.SmallerOperator,
		'=':        token2.SmallerEqualsOperator,
	},
	'%': singleOperatorOption(token2.ModOperator),
	';': singleOperatorOption(token2.SemicolonOperator),
	':': singleOperatorOption(token2.ColonOperator),
	')': singleOperatorOption(token2.RightParenOperator),
	'(': singleOperatorOption(token2.LeftParenOperator),
	'{': singleOperatorOption(token2.LeftCurlyOperator),
	'}': singleOperatorOption(token2.RightCurlyOperator),
	'[': singleOperatorOption(token2.LeftBracketOperator),
	']': singleOperatorOption(token2.RightBracketOperator),
	',': singleOperatorOption(token2.CommaOperator),
	'.': singleOperatorOption(token2.DotOperator),
}

// endOfStatementDisablingOperators are operators that disable the scanners 'insertEos' flag.
// If the scanning gathers one of those operators, it changes the flag to false. The maps
// keys are the disabling operators and their values are the corresponding enabling operators.
var endOfStatementDisablingOperators = map[token2.Operator]token2.Operator{
	token2.LeftParenOperator:   token2.RightParenOperator,
	token2.LeftBracketOperator: token2.RightBracketOperator,
}

// endOfStatementEnablingOperators is a reversed map of the endOfStatementDisablingOperators.
var endOfStatementEnablingOperators map[token2.Operator]token2.Operator

func init() {
	length := len(endOfStatementDisablingOperators)
	endOfStatementEnablingOperators = make(map[token2.Operator]token2.Operator, length)
	for disabler, enabler := range endOfStatementDisablingOperators {
		endOfStatementEnablingOperators[enabler] = disabler
	}
}
