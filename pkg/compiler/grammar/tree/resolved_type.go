package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/scope"

type resolvedType struct {
	symbol *scope.Class
	isResolved bool
}

func (resolvedType *resolvedType) resolve(symbol *scope.Class) {
	resolvedType.symbol = symbol
	resolvedType.isResolved = true
}

func (resolvedType *resolvedType) class() (*scope.Class, bool) {
	return resolvedType.symbol, resolvedType.isResolved
}
