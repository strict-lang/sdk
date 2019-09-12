package backend

import "gitlab.com/strict-lang/sdk/compilation/ast"

type Extension interface {
	ModifyVisitor(generation *Generation, visitor *ast.Visitor)
}

type extensionSet struct {
	elements []Extension
}

func NewExtensionSet(extensions ...Extension) Extension {
	return &extensionSet{elements: extensions}
}

func (extensions *extensionSet) ModifyVisitor(generation *Generation, visitor *ast.Visitor) {
	for _, element := range extensions.elements {
		element.ModifyVisitor(generation, visitor)
	}
}