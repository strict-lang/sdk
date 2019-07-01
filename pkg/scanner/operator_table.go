package scanner

import "github.com/BenjaminNitschke/Strict/pkg/token"

var operatorTable = OperatorTable{
	'+': {
		singleChar: token.AddOperator,
		'=': token.AddAssignOperator,
		'+': token.IncrementOperator,
	},
	'-': {
		singleChar: token.SubOperator,
		'=': token.SubAssignOperator,
		'-': token.DecrementOperator,
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
