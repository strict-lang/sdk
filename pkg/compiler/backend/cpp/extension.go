package cpp

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Extension interface {
	ModifyVisitor(generation *Generation, visitor *tree.DelegatingVisitor)
}

type extensionSet struct {
	elements []Extension
}

func NewExtensionSet(extensions ...Extension) Extension {
	return &extensionSet{elements: extensions}
}

func (extensions *extensionSet) ModifyVisitor(generation *Generation, visitor *tree.DelegatingVisitor) {
	for _, element := range extensions.elements {
		element.ModifyVisitor(generation, visitor)
	}
}
