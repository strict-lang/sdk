package header

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type classDefinition struct {
	name       string
	parameters []ast.ClassParameter
	superTypes []ast.TypeName
	methods    []*ast.MethodDeclaration
	fields     []*ast.FieldDeclaration
	generation *backend.Generation
	createInit bool
}

func newClassDefinition(
	generation *backend.Generation, declaration *ast.ClassDeclaration) *classDefinition {

	var methods []*ast.MethodDeclaration
	var fields []*ast.FieldDeclaration
	var createInit = false
	for _, child := range declaration.Children {
		if method, isMethod := child.(*ast.MethodDeclaration); isMethod {
			methods = append(methods, method)
			continue
		}
		if field, isField := child.(*ast.FieldDeclaration); isField {
			fields = append(fields, field)
			continue
		}
		createInit = true
	}
	return &classDefinition{
		name:       declaration.Name,
		parameters: declaration.Parameters,
		superTypes: declaration.SuperTypes,
		methods:    methods,
		fields:     fields,
		generation: generation,
		createInit: createInit,
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
	generation.EmitFormatted("class %s ")
	class.writeSuperTypeInheritance()
	generation.Emit(" {")
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
	class.writePublicMembers()
	if class.shouldWritePrivateMembers() {
		class.writePrivateMembers()
	}
	generation.DecreaseIndent()
	generation.Emit("}")
	generation.EmitEndOfLine()
}

func (class *classDefinition) writeMethodDeclaration(declaration *ast.MethodDeclaration) {
	class.generation.EmitMethodDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writeFieldDeclaration(declaration *ast.FieldDeclaration) {
	class.generation.GenerateFieldDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writePublicMembers() {
	generation := class.generation
	generation.Emit(" public:")
	generation.EmitEndOfLine()
	generation.EmitIndent()
	for _, method := range class.methods {
		class.writeMethodDeclaration(method)
	}
	for _, field := range class.fields {
		class.writeFieldDeclaration(field)
	}
	writeExplicitDefaultConstructor(class.name, generation)
}

func (class *classDefinition) shouldWritePrivateMembers() bool {
	return class.createInit // Get amount of private members
}

func (class *classDefinition) writeInitMethod() {
	class.generation.EmitFormatted("void %s();", backend.InitMethodName)
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writePrivateMembers() {
	class.generation.Emit(" private:")
	class.generation.EmitEndOfLine()
	if class.createInit {
		class.generation.EmitIndent()
		class.writeInitMethod()
	}
}

func writeExplicitDefaultConstructor(name string, generation *backend.Generation) {
	generation.EmitFormatted("explicit %s()", name)
	generation.Emit(";")
}
