package header

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type ClassDefinition struct {
	name string
	Methods []*ast.MethodDeclaration
	Fields  []*ast.FieldDeclaration
	generation *backend.Generation
}

func (class *ClassDefinition) GenerateCode() {
	generation := class.generation
	generation.EmitFormatted("class %s {", class.name)
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
	class.writePublicMembers()
	class.writePrivateMembers()
	generation.DecreaseIndent()
	generation.Emit("}")
	generation.EmitEndOfLine()
}

func (class *ClassDefinition) writeMethodDeclaration(
	stream backend.OutputStream,
	declaration *ast.MethodDeclaration) {


}

func (class *ClassDefinition) writePublicMembers() {
	generation := class.generation
	generation.Emit(" public:")
	generation.EmitEndOfLine()
	generation.EmitIndent()
	writeExplicitDefaultConstructor(class.name, generation)
}

func (class *ClassDefinition) writePrivateMembers() {
	class.generation.Emit(" private:")
	class.generation.EmitEndOfLine()
}

func writeExplicitDefaultConstructor(name string, generation *backend.Generation) {
	generation.EmitFormatted("explicit %s()", name)
	generation.Emit(";")
}
