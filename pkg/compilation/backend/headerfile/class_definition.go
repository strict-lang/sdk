package headerfile

import (
	backend2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type classDefinition struct {
	name               string
	parameters         []syntaxtree2.ClassParameter
	superTypes         []syntaxtree2.TypeName
	otherMembers       []syntaxtree2.Node
	fields             []syntaxtree2.Node
	generation         *backend2.Generation
	shouldCreateInit   bool
	declarationVisitor *syntaxtree2.Visitor
}

func newClassDefinition(
	generation *backend2.Generation, declaration *syntaxtree2.ClassDeclaration) *classDefinition {

	fields, otherMembers := filterFieldDeclarations(declaration.Children)
	createInit := len(fields) > 0
	definition := &classDefinition{
		name:               declaration.Name,
		parameters:         declaration.Parameters,
		superTypes:         declaration.SuperTypes,
		otherMembers:       otherMembers,
		fields:             fields,
		generation:         generation,
		shouldCreateInit:   createInit,
		declarationVisitor: syntaxtree2.NewEmptyVisitor(),
	}
	initializeVisitor(definition)
	return definition
}

func initializeVisitor(definition *classDefinition) {
	definition.declarationVisitor.VisitMethodDeclaration = definition.writeMethodDeclaration
	definition.declarationVisitor.VisitFieldDeclaration = definition.writeFieldDeclaration
	definition.declarationVisitor.VisitConstructorDeclaration = definition.writeConstructorDeclaration
}

func filterFieldDeclarations(nodes []syntaxtree2.Node) (fields []syntaxtree2.Node, others []syntaxtree2.Node) {
	for _, child := range nodes {
		switch child.(type) {
		case *syntaxtree2.MethodDeclaration, *syntaxtree2.ConstructorDeclaration:
			others = append(others, child)
		case *syntaxtree2.FieldDeclaration:
			fields = append(fields, child)
		}
	}
	return
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
	generation.Emit("};")
	generation.EmitEndOfLine()
}

func (class *classDefinition) writeMethodDeclaration(declaration *syntaxtree2.MethodDeclaration) {
	class.generation.EmitMethodDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writeFieldDeclaration(declaration *syntaxtree2.FieldDeclaration) {
	class.generation.GenerateFieldDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writePublicMembers() {
	generation := class.generation
	generation.Emit(" public:")
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
	writeExplicitDefaultConstructor(class.name, generation)
	for _, member := range class.otherMembers {
		generation.EmitIndent()
		member.Accept(class.declarationVisitor)
	}
	generation.EmitEndOfLine()
	for _, field := range class.fields {
		generation.EmitIndent()
		field.Accept(class.declarationVisitor)
	}
	generation.EmitEndOfLine()
	generation.DecreaseIndent()
}

func (class *classDefinition) shouldWritePrivateMembers() bool {
	return class.shouldCreateInit // Get amount of private members
}

func (class *classDefinition) writeInitMethod() {
	class.generation.EmitFormatted("void %s();", backend2.InitMethodName)
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writePrivateMembers() {
	class.generation.Emit(" private:")
	class.generation.IncreaseIndent()
	class.generation.EmitEndOfLine()
	if class.shouldCreateInit {
		class.generation.EmitIndent()
		class.writeInitMethod()
	}
	class.generation.DecreaseIndent()
}

func (class *classDefinition) writeConstructorDeclaration(declaration *syntaxtree2.ConstructorDeclaration) {
	output := class.generation
	className := output.Unit.Class.Name
	output.Emit("explicit ")
	output.Emit(className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(";")
	output.EmitEndOfLine()
}

func writeExplicitDefaultConstructor(name string, generation *backend2.Generation) {
	generation.EmitFormatted("explicit %s()", name)
	generation.Emit(";")
	generation.EmitEndOfLine()
}