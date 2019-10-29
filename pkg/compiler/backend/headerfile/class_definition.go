package headerfile

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type classDefinition struct {
	name               string
	parameters         []*tree.ClassParameter
	superTypes         []tree.TypeName
	otherMembers       []tree.Node
	fields             []tree.Node
	generation         *backend.Generation
	shouldCreateInit   bool
	declarationVisitor tree.Visitor
}

func newClassDefinition(
	generation *backend.Generation, declaration *tree.ClassDeclaration) *classDefinition {

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
	}
	definition.declarationVisitor = createVisitorForDefinition(definition)
	return definition
}

func createVisitorForDefinition(definition *classDefinition) tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.MethodDeclarationVisitor = definition.writeMethodDeclaration
	visitor.FieldDeclarationVisitor = definition.writeFieldDeclaration
	visitor.ConstructorDeclarationVisitor = definition.writeConstructorDeclaration
	return visitor
}

func filterFieldDeclarations(nodes []tree.Node) (fields []tree.Node, others []tree.Node) {
	for _, child := range nodes {
		switch child.(type) {
		case *tree.MethodDeclaration, *tree.ConstructorDeclaration:
			others = append(others, child)
		case *tree.FieldDeclaration:
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

func (class *classDefinition) writeMethodDeclaration(declaration *tree.MethodDeclaration) {
	class.generation.EmitMethodDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) writeFieldDeclaration(declaration *tree.FieldDeclaration) {
	class.generation.GenerateFieldDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *classDefinition) shouldWriteExplicitDefaultConstructor() bool {
	for _, member := range class.otherMembers {
		if constructor, isConstructor := member.(*tree.ConstructorDeclaration); isConstructor {
			if len(constructor.Parameters) == 0 {
				return false
			}
		}
	}
	return true
}

func (class *classDefinition) writePublicMembers() {
	generation := class.generation
	generation.Emit(" public:")
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
	if class.shouldWriteExplicitDefaultConstructor() {
		writeExplicitDefaultConstructor(class.name, generation)
	}
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
	class.generation.EmitFormatted("void %s();", backend.InitMethodName)
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

func (class *classDefinition) writeConstructorDeclaration(declaration *tree.ConstructorDeclaration) {
	output := class.generation
	className := output.Unit.Class.Name
	output.Emit("explicit ")
	output.Emit(className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(";")
	output.EmitEndOfLine()
}

func writeExplicitDefaultConstructor(name string, generation *backend.Generation) {
	generation.EmitFormatted("explicit %s()", name)
	generation.Emit(";")
	generation.EmitEndOfLine()
}
