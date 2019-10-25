package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type ConcreteTypeName struct {
	Name   string
	Region input.Region
}

func (concrete *ConcreteTypeName) NonGenericName() string {
	return concrete.Name
}

func (concrete *ConcreteTypeName) FullName() string {
	return concrete.Name
}

func (concrete *ConcreteTypeName) Accept(visitor Visitor) {
	visitor.VisitConcreteTypeName(concrete)
}

func (concrete *ConcreteTypeName) AcceptRecursive(visitor Visitor) {
	concrete.Accept(visitor)
}

func (concrete *ConcreteTypeName) Locate() input.Region {
	return concrete.Region
}

func (concrete *ConcreteTypeName) Matches(node Node) bool {
	if target, ok := node.(*ConcreteTypeName); ok {
		return concrete.Name == target. Name
	}
	return false
}
