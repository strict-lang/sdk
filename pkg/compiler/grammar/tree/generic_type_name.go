package tree

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
	"strings"
)

type GenericTypeName struct {
	Name          string
	Arguments     []*Generic
	Region        input.Region
	Parent        Node
	typeReference *TypeReference
}

// TODO: Support Nested Type Names.
// Type names are not expressions but we still want to visit them as such
// we would need an adapter that wraps non expression types and supplies
// functions such as AcceptRecursive and TransformExpression.
type Generic struct {
	Name string
	IsWildcard bool
	Expression Expression
}

const WildcardName = "_wildcard"

func NewWildcardGeneric() *Generic {
	return &Generic{
		Name: WildcardName,
		IsWildcard: true,
		Expression: &WildcardNode{},
	}
}

func NewLetBindingGeneric(binding *LetBinding) *Generic {
	return &Generic{
		Name: binding.Names[0].Value,
		IsWildcard: false,
		Expression: binding,
	}
}

func NewIdentifierGeneric(identifier *Identifier) *Generic {
	return &Generic{
		Name: identifier.Value,
		IsWildcard: false,
		Expression: identifier,
	}
}

func (name *GenericTypeName) TypeReference() *TypeReference {
	return name.typeReference
}

func (name *GenericTypeName) SetEnclosingNode(target Node) {
	name.Parent = target
}

func (name *GenericTypeName) EnclosingNode() (Node, bool) {
	return name.Parent, name.Parent != nil
}

func (name *GenericTypeName) FullName() string {
	return fmt.Sprintf("%s<%s>", name.Name, name.joinArguments())
}

func (name *GenericTypeName) joinArguments() string {
	if len(name.Arguments) == 0 {
		return ""
	}
	var builder strings.Builder
	builder.WriteString(name.Arguments[0].Name)
	for _, argument := range name.Arguments[1:]{
		builder.WriteString(", ")
		builder.WriteString(argument.Name)
	}
	return builder.String()
}

func (name *GenericTypeName) BaseName() string {
	return name.Name
}

func (name *GenericTypeName) Accept(visitor Visitor) {
	visitor.VisitGenericTypeName(name)
}

func (name *GenericTypeName) AcceptRecursive(visitor Visitor) {
	name.Accept(visitor)
	for _, argument := range name.Arguments {
		argument.Expression.AcceptRecursive(visitor)
	}
}

func (name *GenericTypeName) Locate() input.Region {
	return name.Region
}

func (name *GenericTypeName) Matches(node Node) bool {
	if target, ok := node.(*GenericTypeName); ok {
		return name.Name == target.Name && name.matchesArguments(target.Arguments)
	}
	return false
}

func (name *GenericTypeName) matchesArguments(arguments []*Generic) bool {
	if len(arguments) != len(name.Arguments) {
		return false
	}
	for index, argument := range name.Arguments {
		if !argument.Expression.Matches(arguments[index].Expression) {
			return false
		}
	}
	return true
}
