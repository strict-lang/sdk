package analysis

import "strict.dev/sdk/pkg/compiler/scope"

type Importing interface {
	Import(scope scope.MutableScope) error
}
