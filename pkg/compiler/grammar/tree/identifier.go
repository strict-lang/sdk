package tree

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

type Identifier struct {
	Value        string
	Region       input.Region
	resolvedType resolvedType
	Parent Node
	binding scope.Symbol
	inDeclaration bool
}

func (identifier *Identifier) ReferencePoint() scope.ReferencePoint {
	return scope.NewReferencePointWithPosition(
		identifier.Value, identifier.Region.Begin())
}

func (identifier *Identifier) IsPartOfDeclaration() bool {
	return identifier.inDeclaration
}

func (identifier *Identifier) MarkAsPartOfDeclaration() {
	identifier.inDeclaration = true
}

func (identifier *Identifier) Bind(target scope.Symbol) {
	identifier.binding = target
}

func (identifier *Identifier) Binding() scope.Symbol {
	return identifier.binding
}

func (identifier *Identifier) IsBound() bool {
	return identifier.binding != nil
}

func (identifier *Identifier) SetEnclosingNode(target Node) {
  identifier.Parent = target
}

func (identifier *Identifier) EnclosingNode() (Node, bool) {
  return identifier.Parent, identifier.Parent != nil
}

func (identifier *Identifier) Accept(visitor Visitor) {
	visitor.VisitIdentifier(identifier)
}

func (identifier *Identifier) AcceptRecursive(visitor Visitor) {
	identifier.Accept(visitor)
}

func (identifier *Identifier) Locate() input.Region {
	return identifier.Region
}

func (identifier *Identifier) Matches(node Node) bool {
	if target, ok := node.(*Identifier); ok {
		return identifier.Value == target.Value
	}
	return false
}
