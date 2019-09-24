package backend

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Extension interface {
	ModifyVisitor(generation *Generation, visitor *syntaxtree2.Visitor)
}

type extensionSet struct {
	elements []Extension
}

func NewExtensionSet(extensions ...Extension) Extension {
	return &extensionSet{elements: extensions}
}

func (extensions *extensionSet) ModifyVisitor(generation *Generation, visitor *syntaxtree2.Visitor) {
	for _, element := range extensions.elements {
		element.ModifyVisitor(generation, visitor)
	}
}
