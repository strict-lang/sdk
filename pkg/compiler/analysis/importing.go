package analysis

import "github.com/strict-lang/sdk/pkg/compiler/scope"

type Importing interface {
	Import(scope scope.MutableScope) error
}
