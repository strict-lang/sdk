package tree

import "github.com/strict-lang/sdk/pkg/compiler/input"

type FieldDeclaration struct {
	Name     *Identifier
	TypeName TypeName
	Region   input.Region
	Parent   Node
	Inferred bool
}

func (field *FieldDeclaration) SetEnclosingNode(target Node) {
	field.Parent = target
}

func (field *FieldDeclaration) EnclosingNode() (Node, bool) {
	return field.Parent, field.Parent != nil
}

func (field *FieldDeclaration) Accept(visitor Visitor) {
	visitor.VisitFieldDeclaration(field)
}

func (field *FieldDeclaration) AcceptRecursive(visitor Visitor) {
	field.Accept(visitor)
	field.Name.AcceptRecursive(visitor)
	field.TypeName.AcceptRecursive(visitor)
}

func (field *FieldDeclaration) Locate() input.Region {
	return field.Region
}

func (field *FieldDeclaration) Matches(node Node) bool {
	if target, ok := node.(*FieldDeclaration); ok {
		return field.Name.Matches(target.Name) &&
			field.TypeName.Matches(target.TypeName)
	}
	return false
}
