package tree

import (
	"strict.dev/sdk/pkg/compiler/input"
	"strict.dev/sdk/pkg/compiler/typing"
)

type TypeName interface {
	Node
	// FullName returns the full type name, including generics and punctuation.
	FullName() string
	// BaseName returns base type names, without punctuation and generics.
	// Example: Number[] -> Number, MutableList<String> -> MutableList.
	BaseName() string
	TypeReference() *TypeReference
}

func ParseTypeName(region input.Region, deducedType typing.Type) TypeName {
	parser := &typeNameParser{
		region: region,
	}
	return parser.parse(deducedType)
}

type typeNameParser struct {
	region input.Region
	lastType TypeName
}

func (parser *typeNameParser) parse(value typing.Type) TypeName {
	currentLast := parser.lastType
	value.Accept(parser)
	parsed := parser.lastType
	parser.lastType = currentLast
	return parsed
}

func (parser *typeNameParser) VisitOptional(optional *typing.OptionalType) {
	child := parser.parse(optional.Child)
	parser.lastType = &OptionalTypeName{
		Region:        parser.region,
		TypeName:      child,
		typeReference: &TypeReference{resolved: optional},
	}
}

func (parser *typeNameParser) VisitConcrete(concrete *typing.ConcreteType) {
	parser.lastType = &ConcreteTypeName{
		Region:        parser.region,
		Name:     concrete.Name,
		typeReference: &TypeReference{resolved: concrete},
	}
}

func (parser *typeNameParser) VisitGeneric(generic *typing.GenericType) {
	generics := make([]TypeName, len(generic.Arguments))
	for index, argument := range generic.Arguments {
		generics[index] = parser.parse(argument)
	}
	parser.lastType = &GenericTypeName{
		Name:          generic.Concrete().String(),
		Generic:       generics[0],
		Region:        parser.region,
		typeReference: &TypeReference{resolved: generic},
	}
}

func (parser *typeNameParser) VisitList(list *typing.ListType) {
	parser.lastType = &ListTypeName{
		Element:       parser.parse(list.Child),
		Region:        parser.region,
		typeReference: &TypeReference{resolved: list},
	}
}

type TypeReference struct {
	resolved typing.Type
}

func (reference *TypeReference) Resolved() (typing.Type, bool) {
	return reference.resolved, reference.resolved != nil
}

func (reference *TypeReference) Resolve(target typing.Type) {
	reference.resolved = target
}
