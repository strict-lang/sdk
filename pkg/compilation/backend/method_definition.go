package backend

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	"strings"
)

type MethodDefinition struct {
	generation  *Generation
	declaration *syntaxtree.MethodDeclaration
	buffer      *strings.Builder
	prologue    map[string]ProloguePart
	epilogue    map[string]EpiloguePart
}

type ProloguePart func()
type EpiloguePart func()

func (generation *Generation) NewMethodGeneration(method *syntaxtree.MethodDeclaration) *MethodDefinition {
	return &MethodDefinition{
		generation:  generation,
		declaration: method,
		buffer:      &strings.Builder{},
		epilogue:    map[string]EpiloguePart{},
		prologue:    map[string]ProloguePart{},
	}
}

func (generation *Generation) GenerateMethod(method *syntaxtree.MethodDeclaration) {
	methodGenerator := generation.NewMethodGeneration(method)
	generation.method = methodGenerator
	methodGenerator.Emit()
	generation.method = nil
}

func (definition *MethodDefinition) Emit() {
	generator := definition.generation
	generator.buffer = definition.buffer

	declaration := definition.generateDeclaration()
	generator.IncreaseIndent()
	methodBody := definition.generateBody()
	prologue := definition.generatePrologue()
	epilogue := definition.generateEpilogue()

	generator.buffer = generator.output
	generator.Emit(declaration)
	generator.Emit(" {\n")
	generator.EmitIndent()
	generator.Emit(prologue)
	generator.Emit(methodBody)
	generator.Emit(epilogue)
	generator.DecreaseIndent()
	generator.EmitIndent()
	generator.Emit("}\n")
}

func (definition *MethodDefinition) generateDeclaration() string {
	definition.buffer.Reset()
	generation := definition.generation
	generation.EmitMethodDeclaration(definition.declaration)
	return definition.buffer.String()
}

func (definition *MethodDefinition) generateBody() string {
	definition.buffer.Reset()
	// Do not use the normal BlockStatement visitor for BlockStatement definition in
	// a method body. It will generate open and close brackets and this will produce
	// faulty code, if a prologue or epilogue is generated.
	if block, ok := definition.declaration.Body.(*syntaxtree.BlockStatement); ok {
		for _, child := range block.Children {
			definition.generation.EmitNode(child)
		}
		return definition.buffer.String()
	}
	definition.generation.EmitNode(definition.declaration.Body)
	return definition.buffer.String()
}

func (definition *MethodDefinition) generatePrologue() string {
	definition.buffer.Reset()
	for _, part := range definition.prologue {
		part()
	}
	if len(definition.prologue) > 0 {
		definition.generation.Emit("\n")
		definition.generation.EmitIndent()
	}
	return definition.buffer.String()
}

func (definition *MethodDefinition) generateEpilogue() string {
	definition.buffer.Reset()
	for _, part := range definition.epilogue {
		part()
	}
	return definition.buffer.String()
}

func (definition *MethodDefinition) addToEpilogue(name string, function EpiloguePart) {
	definition.epilogue[name] = function
}

func (definition *MethodDefinition) addToPrologue(name string, function ProloguePart) {
	definition.prologue[name] = function
}
