package scanner

import "github.com/BenjaminNitschke/Strict/pkg/token"

var complexOperatorScanners = map[rune]map[rune] token.Kind {
	'+': {
		'=': token.Assign,
	},
}

func (scanner Scanner) scanOperator() {

}
