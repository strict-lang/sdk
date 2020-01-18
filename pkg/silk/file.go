package silk

import "gitlab.com/strict-lang/sdk/pkg/silk/symbol"

type Version int

type File struct {
	Symbols *symbol.Table
	ClassDefinition *ClassDefinition
	Version Version
}