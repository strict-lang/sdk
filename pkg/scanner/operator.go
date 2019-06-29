package scanner

import "github.com/BenjaminNitschke/Strict/pkg/token"

var complexOperatorScanners = map[rune]map[rune]token.Operator{
	'+': {
		'=': token.AssignOperator,
	},
}

func (scanner Scanner) scanOperator() {

}
