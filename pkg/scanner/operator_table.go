package scanner

import "github.com/BenjaminNitschke/Strict/pkg/token"

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
		'>':        token.ShiftRightOperator,
	},
	'<': {
		singleChar: token.SmallerOperator,
		'=':        token.SmallerEqualsOperator,
		'<':        token.ShiftLeftOperator,
	},
	';': singleOperatorOption(token.SemicolonOperator),
	':': singleOperatorOption(token.ColonOperator),
	'(': singleOperatorOption(token.LeftParenOperator),
	')': singleOperatorOption(token.RightParenOperator),
	'{': singleOperatorOption(token.LeftCurlyOperator),
	'}': singleOperatorOption(token.RightCurlyOperator),
	'[': singleOperatorOption(token.LeftBracketOperator),
	']': singleOperatorOption(token.RightBracketOperator),
}
