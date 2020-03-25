package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/scope"

type resolvedType struct {
	symbol *scope.Class
}

func (resolvedType *resolvedType) resolve(symbol *scope.Class) {
	resolvedType.symbol = symbol
}

func (resolvedType *resolvedType) class() (*scope.Class, bool) {
	return resolvedType.symbol, resolvedType.symbol != nil
}
