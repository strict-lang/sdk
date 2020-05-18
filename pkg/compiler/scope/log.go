package scope

import "log"

func Log(scope Scope) {
	symbols := scope.Search(func(symbol Symbol) bool {
		return true
	})
	for _, symbol := range symbols {
		log.Printf("- %s:%s = %+v", symbol.scopeId, symbol.Symbol.Name(), symbol.Symbol)
	}
}
