package typing

import "fmt"

type ListType struct {
	Child Type
}

func (list *ListType) Concrete() Type {
	return list.Child.Concrete()
}

func (list *ListType) String() string {
	return fmt.Sprintf("%s[]", list.Child)
}

func (list *ListType) Is(target Type) bool {
	if targetList, ok := target.(*ListType); ok {
		return list.Child.Is(targetList.Child)
	}
	return false
}

func (list *ListType) Accept(visitor Visitor) {
	visitor.VisitList(list)
}

func (list *ListType) AcceptRecursive(visitor Visitor) {
	list.Accept(visitor)
	list.Child.AcceptRecursive(visitor)
}

