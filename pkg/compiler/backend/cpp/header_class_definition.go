package cpp

import "github.com/strict-lang/sdk/pkg/compiler/grammar/tree"

type headerClass struct {
	name               string
	parameters         []*tree.ClassParameter
	superTypes         []tree.TypeName
	otherMembers       []tree.Node
	fields             []tree.Node
	generation         *Generation
	shouldCreateInit   bool
	trait              bool
	declarationVisitor tree.Visitor
}

func newClassDefinition(
	generation *Generation, declaration *tree.ClassDeclaration) *headerClass {

	fields, otherMembers := filterFieldDeclarations(declaration.Children)
	createInit := !declaration.Trait
	definition := &headerClass{
		name:             declaration.Name,
		parameters:       declaration.Parameters,
		superTypes:       declaration.SuperTypes,
		otherMembers:     otherMembers,
		fields:           fields,
		trait:            declaration.Trait,
		generation:       generation,
		shouldCreateInit: createInit,
	}
	definition.declarationVisitor = createVisitorForDefinition(definition)
	return definition
}

func createVisitorForDefinition(definition *headerClass) tree.Visitor {
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

func (class *headerClass) writeTemplates() {
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

func (class *headerClass) writeSuperTypeInheritance() {
	for index, superType := range class.superTypes {
		if index != 0 {
			class.generation.Emit(", ")
		}
		class.generation.EmitNode(superType)
	}
}

func (class *headerClass) generateCode() {
	generation := class.generation
	generation.EmitFormatted("class %s ", class.name)
	if !class.trait {
		generation.Emit("final ")
	}
	class.writeSuperTypeInheritance()
	generation.Emit("{")
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	class.writePublicMembers()
	class.writePrivateMembers()
	generation.DecreaseIndent()
	generation.Emit("};")
	generation.EmitEndOfLine()
}

func (class *headerClass) writeMethodDeclaration(declaration *tree.MethodDeclaration) {
	class.generation.EmitMethodDeclaration(declaration)
	if class.trait {
		class.generation.Emit(" = 0")
	}
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *headerClass) writeFieldDeclaration(declaration *tree.FieldDeclaration) {
	class.generation.GenerateFieldDeclaration(declaration)
	class.generation.Emit(";")
	class.generation.EmitEndOfLine()
}

func (class *headerClass) shouldWriteExplicitDefaultConstructor() bool {
	if class.trait {
		return false
	}
	for _, member := range class.otherMembers {
		if constructor, isConstructor := member.(*tree.ConstructorDeclaration); isConstructor {
			if len(constructor.Parameters) == 0 {
				return false
			}
		}
	}
	return true
}

func (class *headerClass) writePublicMembers() {
	generation := class.generation
	generation.Emit("public:")
	generation.EmitEndOfLine()
	if class.shouldWriteExplicitDefaultConstructor() {
		generation.EmitIndent()
		writeExplicitDefaultConstructor(class.name, generation)
	}
	for _, member := range class.otherMembers {
		member.Accept(class.declarationVisitor)
	}
	if len(class.otherMembers) != 0 && len(class.fields) != 0 {
		generation.EmitEndOfLine()
	}
	for _, field := range class.fields {
		field.Accept(class.declarationVisitor)
	}
}

func (class *headerClass) shouldWritePrivateMembers() bool {
	return class.shouldCreateInit // Get amount of private members
}

func (class *headerClass) writeInitMethod() {
	class.generation.EmitFormatted("void %s();", InitMethodName)
	class.generation.EmitEndOfLine()
}

func (class *headerClass) writePrivateMembers() {
	if class.trait {
		return
	}
	class.generation.Emit("private:")
	class.generation.EmitEndOfLine()
	class.generation.EmitIndent()
	class.writeInitMethod()
}

func (class *headerClass) writeConstructorDeclaration(declaration *tree.ConstructorDeclaration) {
	output := class.generation
	className := output.Unit.Class.Name
	output.Emit("explicit ")
	output.Emit(className)
	output.EmitParameterList(declaration.Parameters)
	output.Emit(";")
	output.EmitEndOfLine()
}

func writeExplicitDefaultConstructor(name string, generation *Generation) {
	generation.EmitFormatted("explicit %s()", name)
	generation.Emit(";")
	generation.EmitEndOfLine()
	generation.EmitIndent()
	generation.EmitEndOfLine()
}
