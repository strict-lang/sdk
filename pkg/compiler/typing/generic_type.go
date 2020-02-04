package typing

import (
	"fmt"
	"strings"
)

type GenericType struct {
	Child     Type
	Arguments []Type
}

func (generic *GenericType) Is(target Type) bool {
	if targetGeneric, ok := target.(*GenericType); ok {
		return generic.matches(targetGeneric)
	}
	return false
}

func (generic *GenericType) matches(target *GenericType) bool {
	return generic.Child.Is(target.Child) && generic.matchesArguments(target)
}

func (generic *GenericType) matchesArguments(target *GenericType) bool {
	if len(generic.Arguments) != len(target.Arguments) {
		return false
	}
	return generic.sharesTraits(target)
}

func (generic *GenericType) sharesTraits(target *GenericType) bool {
	for index, argument := range target.Arguments {
		if targetArgument := target.Arguments[index]; targetArgument != argument {
			return false
		}
	}
	return true
}

func (generic *GenericType) Accept(visitor Visitor) {
	visitor.VisitGeneric(generic)
}

func (generic *GenericType) AcceptRecursive(visitor Visitor) {
	generic.Accept(visitor)
	generic.Child.AcceptRecursive(visitor)
	for _, argument := range generic.Arguments {
		argument.AcceptRecursive(visitor)
	}
}

func (generic *GenericType) Concrete() Type {
	return generic.Child.Concrete()
}

func (generic *GenericType) String() string {
	argumentNames := make([]string, len(generic.Arguments))
	for index, argument := range generic.Arguments {
		argumentNames[index] = argument.String()
	}
	formatted := strings.Join(argumentNames, ", ")
	return fmt.Sprintf("%s<%s>", generic.Child, formatted)
}
