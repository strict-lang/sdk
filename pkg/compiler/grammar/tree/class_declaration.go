package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ClassDeclaration struct {
	Name       string
	Parameters []ClassParameter
	SuperTypes []TypeName
	Children   []Node
	Region     input.Region
}

type ClassParameter struct {
	Name      string
	SuperType TypeName
}

func (class *ClassDeclaration) Accept(visitor Visitor) {
	visitor.VisitClassDeclaration(class)
}

func (class *ClassDeclaration) AcceptRecursive(visitor Visitor) {
	class.Accept(visitor)
	for _, child := range class.Children {
		child.AcceptRecursive(visitor)
	}
}

func (class *ClassDeclaration) Locate() input.Region {
	return class.Region
}
