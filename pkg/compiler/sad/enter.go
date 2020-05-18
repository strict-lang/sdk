package sad

import (
	"github.com/strict-lang/sdk/pkg/compiler/scope"
	"github.com/strict-lang/sdk/pkg/compiler/typing"
	"strings"
)

type ClassFilter func(*Class) bool

var noFilter ClassFilter = func(class *Class) bool {
	return true
}

// Enters the members that have been defined in aan api-descriptor file
// into a given scope. Used during importing of foreign namespaces.
func Enter(tree *Tree, scope scope.MutableScope)  {
	EnterFiltered(tree, scope, noFilter)
}

func EnterFiltered(tree *Tree, scope scope.MutableScope, filter ClassFilter) {
	entering := newEntering(tree, scope, filter)
	entering.enter()
}

type entering struct {
	scope scope.MutableScope
	tree  *Tree
	filter ClassFilter
	classes map[string] classBinding
}

type classBinding struct {
	class *Class
	symbol *scope.Class
}

func newEntering(tree *Tree, scope scope.MutableScope, filter ClassFilter) *entering {
	return &entering{
		scope:   scope,
		tree:    tree,
		filter:  filter,
		classes: map[string]classBinding{},
	}
}

func (entering *entering) enter() {
	entering.enterClasses()
	entering.populateClasses()
}

func (entering *entering) createMethodSymbol(
	enclosingClass *scope.Class,
	method Method) *scope.Method {

	return &scope.Method{
		DeclarationName: method.Name,
		ReturnType:      entering.findClass(method.ReturnType),
		Parameters:      entering.translateParameters(enclosingClass, method.Parameters),
	}
}

func (entering *entering) translateParameters(
	enclosingClass *scope.Class,
	parameters []Parameter) (fields []*scope.Field) {

	for _, parameter := range parameters {
		field := &scope.Field{
			DeclarationName: parameter.Name,
			Class:           entering.findClass(parameter.Class),
			Kind:            scope.ParameterField,
			EnclosingClass:  enclosingClass,
		}
		fields = append(fields, field)
	}
	return
}

func (entering *entering) createFieldSymbol(class *scope.Class, field Field) *scope.Field {
	return &scope.Field{
		DeclarationName: field.Name,
		Class:           entering.findClass(field.Class),
		Kind:            scope.MemberField,
		EnclosingClass:  class,
	}
}

func (entering *entering) enterClasses() {
	for _, class := range entering.tree.Classes {
		if entering.filter(class) {
			entering.enterClass(class)
		}
	}
}

func (entering *entering) enterClass(class *Class) {
	symbol := &scope.Class{
		Scope:           entering.scope,
		DeclarationName: removeQualifier(class.Name),
		ActualClass:     entering.createTypeInformation(class),
	}
	entering.scope.Insert(symbol)
	entering.classes[class.Name] = classBinding{class:  class, symbol: symbol}
}

func (entering *entering) createTypeInformation(class *Class) typing.Type {
	return &typing.ConcreteType{
		Name:   class.Name,
	}
}

func removeQualifier(name string) string {
	qualifierEnd := strings.Index(name, ".")
	if qualifierEnd == -1 {
		return name
	}
	return name[qualifierEnd:]
}

func (entering *entering) populateClasses() {
	for _, binding := range entering.classes {
		entering.populateClassSymbol(binding)
	}
}

func (entering *entering) populateClassSymbol(binding classBinding) {
	if concrete, ok := binding.symbol.ActualClass.(*typing.ConcreteType); ok {
		concrete.Traits = entering.findActualClasses(binding.class.Traits)
	}
	entering.insertMembers(binding)
}

func (entering *entering) insertMembers(binding classBinding) {
	classScope := binding.symbol.Scope
	for _, method := range binding.class.Methods {
		classScope.Insert(entering.createMethodSymbol(binding.symbol, method))
	}
	for _, field := range binding.class.Fields {
		classScope.Insert(entering.createFieldSymbol(binding.symbol, field))
	}
}

// TODO: Implement generic types
func (entering *entering) findActualClasses(names []ClassName) (classes []typing.Type) {
	for _, name := range names {
		class := entering.findClassByName(name.Name)
		classes = append(classes, class.ActualClass)
	}
	return
}

func (entering *entering) findClass(name ClassName) *scope.Class {
	return entering.findClassByName(name.Name)
}

func (entering *entering) findClassByName(name string) *scope.Class {
	symbols := entering.scope.Lookup(scope.NewReferencePoint(name))
	if !symbols.IsEmpty() {
		if class, ok := scope.AsClassSymbol(symbols.First().Symbol); ok {
			return class
		}
	}
	return scope.Builtins.Any
}
