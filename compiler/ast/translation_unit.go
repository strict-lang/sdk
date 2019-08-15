package ast

import (
	"gitlab.com/strict-lang/sdk/compiler/code"
)

// TranslationUnit represents a unit of translation, a file containing strict
// source code. It can have multiple children, which are seen as the roots
// of the actual ast. This node however, is the real unit of the ast.
type TranslationUnit struct {
	name     string
	scope    *code.Scope
	Children []Node
	NodePosition Position
}

func NewEmptyTranslationUnit(name string) *TranslationUnit {
	return NewTranslationUnit(name, code.NewRootScope(), []Node{})
}

func NewTranslationUnit(name string, rootScope *code.Scope, children []Node) *TranslationUnit {
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

func (unit *TranslationUnit) Scope() *code.Scope {
	return unit.scope
}

func (unit *TranslationUnit) Accept(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
}

func (unit *TranslationUnit) AcceptAll(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
	for _, topLevelNode := range unit.Children {
		topLevelNode.AcceptAll(visitor)
	}
}

func (unit *TranslationUnit) AppendChild(node Node) {
	unit.Children = append(unit.Children, node)
}

func (unit *TranslationUnit) ToTypeName() TypeName {
	return ConcreteTypeName{
		Name: unit.name,
	}
}

func (unit *TranslationUnit) Position() Position {
	return unit.NodePosition
}