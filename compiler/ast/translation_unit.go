package ast

import (
	"gitlab.com/strict-lang/sdk/compiler/scope"
)

// TranslationUnit represents a unit of translation, a file containing strict
// source code. It can have multiple children, which are seen as the roots
// of the actual ast. This node however, is the real unit of the ast.
type TranslationUnit struct {
	name     string
	scope    *scope.Scope
	Children []Node
}

func NewEmptyTranslationUnit(name string) *TranslationUnit {
	return NewTranslationUnit(name, scope.NewRoot(), []Node{})
}

func NewTranslationUnit(name string, rootScope *scope.Scope, children []Node) *TranslationUnit {
	childScope := rootScope.NewNamedChild(name)
	return &TranslationUnit{
		name:     name,
		scope:    childScope,
		Children: children,
	}
}

func (unit *TranslationUnit) Name() string {
	return unit.name
}

func (unit *TranslationUnit) Scope() *scope.Scope {
	return unit.scope
}

func (unit *TranslationUnit) Accept(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
}

func (unit *TranslationUnit) AppendChild(node Node) {
	unit.Children = append(unit.Children, node)
}
