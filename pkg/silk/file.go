package silk

import "strict.dev/sdk/pkg/silk/symbol"

type Version int

type File struct {
	Symbols *symbol.Table
	Class   *Class
	Version Version
}
