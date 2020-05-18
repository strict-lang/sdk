package sad

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"log"
)

func Generate(unit *tree.TranslationUnit) *Class {
	return newGeneration(unit).Generate()
}

type generation struct {
	visitor tree.Visitor
	unit    *tree.TranslationUnit
	class   *Class
}

func newGeneration(unit *tree.TranslationUnit) *generation {
	generation := &generation{unit: unit}
	generation.visitor = generation.createVisitor()
	generation.class = generation.createClass()
	return generation
}

func (generation *generation) Generate() *Class {
	generation.unit.AcceptRecursive(generation.visitor)
	return generation.class
}

func (generation *generation) createVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.MethodDeclarationVisitor = generation.visitMethod
	visitor.FieldDeclarationVisitor = generation.visitField
	return visitor
}

func (generation *generation) createClass() *Class {
	return &Class{
		Kind:    translateKind(generation.unit),
		Traits:  translateTypeNames(generation.unit.Class.SuperTypes),
		Name:    generation.unit.Class.Name,
		Methods: map[string]Method{},
		Fields:  map[string]Field{},
	}
}

func translateTypeNames(names []tree.TypeName) (output []ClassName) {
	for _, name := range names {
		output = append(output, translateTypeName(name))
	}
	return
}

func translateKind(unit *tree.TranslationUnit) TypeKind {
	if unit.Class.Trait {
		return TraitKind
	}
	return ClassKind
}

func (generation *generation) visitMethod(method *tree.MethodDeclaration) {
	descriptor := Method{
		Name:       method.Name.Value,
		Parameters: translateParameters(method),
		ReturnType: translateTypeName(method.Type),
	}
	generation.class.Methods[descriptor.Name] = descriptor
}

func (generation *generation) visitField(field *tree.FieldDeclaration) {
	descriptor := Field{
		Name:  field.Name.Value,
		Class: translateTypeName(field.TypeName),
	}
	generation.class.Fields[descriptor.Name] = descriptor
}

func translateParameters(method *tree.MethodDeclaration) (parameters []Parameter) {
	for _, parameter := range method.Parameters {
		parameters = append(parameters, Parameter{
			Name:  parameter.Name.Value,
			Class: translateTypeName(parameter.Type),
		})
	}
	return
}

var voidType = ClassName{Name: "Strict.Base.Void"}
const sliceTypeName = "Strict.Base.Slice.Slice"

func translateTypeName(name tree.TypeName) ClassName {
	if name == nil {
		return voidType
	}
	switch concrete := name.(type) {
	case *tree.ConcreteTypeName:
		return ClassName{Name: concrete.FullName()}
	case *tree.GenericTypeName:
		return ClassName{
			Name: concrete.FullName(),
			Arguments: translateGenerics(concrete.Arguments),
		}
	case *tree.ListTypeName:
		return ClassName{
			Name:      sliceTypeName,
			Arguments: []ClassName{translateTypeName(concrete.Element)},
		}
	default:
		log.Fatalf("invalid type during sad-generation: %v", concrete)
		return ClassName{}
	}
}

func translateGenerics(generics []*tree.Generic) (names []ClassName) {
	for _, generic := range generics {
		names = append(names, translateGeneric(generic))
	}
	return
}

func translateGeneric(generic *tree.Generic) ClassName {
	if generic.IsWildcard {
		return ClassName{Wildcard: true}
	}
	return ClassName{Name: generic.Name}
}
