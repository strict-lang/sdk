package ast

// TranslationUnit represents a unit of translation, a file containing strict
// source code. It can have multiple children, which are seen as the roots
// of the actual ast. This node however, is the real unit of the ast.
type TranslationUnit struct {
	name     string
	scope    *Scope
	span     Position
	Children []Node
}

func NewTranslationUnit(name string, rootScope Scope, span Position, children []Node) TranslationUnit {
	scope := rootScope.NewNamedChild(name)
	return TranslationUnit{
		name:     name,
		scope:    scope,
		span:     span,
		Children: children,
	}
}

func (unit *TranslationUnit) Name() string {
	return unit.name
}

func (unit *TranslationUnit) Scope() *Scope {
	return unit.scope
}

func (unit *TranslationUnit) Accept(visitor *Visitor) {
	visitor.VisitTranslationUnit(unit)
	for _, child := range unit.Children {
		child.Accept(visitor)
	}
}
