package headerfile

import (
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type classDefinition struct {
	name             string
	parameters       []syntaxtree.ClassParameter
	superTypes       []syntaxtree.TypeName
	methods          []*syntaxtree.MethodDeclaration
	fields           []*syntaxtree.FieldDeclaration
	generation       *backend.Generation
	shouldCreateInit bool
}

func newClassDefinition(
	generation *backend.Generation, declaration *syntaxtree.ClassDeclaration) *classDefinition {

	var methods []*syntaxtree.MethodDeclaration
	var fields []*syntaxtree.FieldDeclaration
	var createInit = false
	for _, child := range declaration.Children {
		if method, isMethod := child.(*syntaxtree.MethodDeclaration); isMethod {
			methods = append(methods, method)
			continue
		}
		if field, isField := child.(*syntaxtree.FieldDeclaration); isField {
			fields = append(fields, field)
			continue
		}
		createInit = true
	}
	return &classDefinition{
		name:             declaration.Name,
		parameters:       declaration.Parameters,
		superTypes:       declaration.SuperTypes,
		methods:          methods,
		fields:           fields,
		generation:       generation,
		shouldCreateInit: createInit,
	}
}

func (class *classDefinition) writeTemplates() {
	if len(class.parameters) == 0 {
		return
	}
	class.generation.Emit("template <")
	for index, parameter := range class.parameters {
		if index != 0 {
			class.generation.Emit(", ")
		}
		class.generation.EmitFormatted("typename %s", parameter.Name)
	}
	class.generation.Emit(">")
}

func (class *classDefinition) writeSuperTypeInheritance() {
	for index, superType := range class.superTypes {
		if index != 0 {
			class.generation.Emit(", ")
		}
		class.generation.EmitNode(superType)
	}
}

func (class *classDefinition) generateCode() {
	generation := class.generation
	generation.EmitFormatted("class %s ", class.name)
	class.writeSuperTypeInheritance()
	generation.Emit("{")
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	class.writePublicMembers()
	if class.shouldWritePrivateMembers() {
		class.writePrivateMembers()
	}
	generation.DecreaseIndent()
	generation.Emit("}")
	generation.EmitEndOfLine()
}

func (class *classDefinition) writeMethodDeclaration(declaration *syntaxtree.MethodDeclaration) {
	class.generation.EmitMethodDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writeFieldDeclaration(declaration *syntaxtree.FieldDeclaration) {
	class.generation.GenerateFieldDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writePublicMembers() {
	generation := class.generation
	generation.Emit(" public:")
	generation.EmitEndOfLine()
	for _, method := range class.methods {
		generation.EmitIndent()
		class.writeMethodDeclaration(method)
	}
	for _, field := range class.fields {
		generation.EmitIndent()
		class.writeFieldDeclaration(field)
	}
	generation.EmitIndent()
	writeExplicitDefaultConstructor(class.name, generation)
	generation.EmitEndOfLine()
}

func (class *classDefinition) shouldWritePrivateMembers() bool {
	return class.shouldCreateInit // Get amount of private members
}

func (class *classDefinition) writeInitMethod() {
	class.generation.EmitFormatted("void %s();", backend.InitMethodName)
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writePrivateMembers() {
	class.generation.Emit(" private:")
	class.generation.EmitEndOfLine()
	if class.shouldCreateInit {
		class.generation.EmitIndent()
		class.writeInitMethod()
	}
}

func writeExplicitDefaultConstructor(name string, generation *backend.Generation) {
	generation.EmitFormatted("explicit %s()", name)
	generation.Emit(";")

}
