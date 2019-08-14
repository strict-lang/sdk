package codegen

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"strings"
)

type MethodGeneration struct {
	generator          *CodeGenerator
	prologueGenerators map[string]PrologueGenerator
	epilogueGenerators map[string]EpilogueGenerator
	declaration        ast.Method
	buffer             *strings.Builder
}

type PrologueGenerator func()
type EpilogueGenerator func()

func (generator *CodeGenerator) NewMethodGeneration(method ast.Method) *MethodGeneration {
	return &MethodGeneration{
		generator:          generator,
		declaration:        method,
		buffer:             &strings.Builder{},
		epilogueGenerators: map[string]EpilogueGenerator{},
		prologueGenerators: map[string]PrologueGenerator{},
	}
}

func (generator *CodeGenerator) GenerateMethod(method *ast.Method) {
	methodGenerator := generator.NewMethodGeneration(*method)
	generator.method = methodGenerator
	methodGenerator.Complete()
	generator.method = nil
}

func (generation *MethodGeneration) Complete() {
	generator := generation.generator
	generator.buffer = generation.buffer

	declaration := generation.generateDeclaration()
	generator.enterBlock()
	methodBody := generation.generateBody()
	prologue := generation.generatePrologue()
	epilogue := generation.generateEpilogue()

	generator.buffer = generator.output
	generator.Emit(declaration)
	generator.Emit(" {\n")
	generator.Spaces()
	generator.Emit(prologue)
	generator.Emit(methodBody)
	generator.Emit(epilogue)
	generator.leaveBlock()
	generator.Spaces()
	generator.Emit("}\n")
}

func (generation *MethodGeneration) generateDeclaration() string {
	generation.buffer.Reset()
	generator := generation.generator

	method := generation.declaration
	returnTypeName := updateTypeName(method.Type)

	generator.Spaces()
	generator.Emitf("%s %s(", returnTypeName.FullName(), method.Name.Value)
	for index, parameter := range method.Parameters {
		if index != 0 {
			generator.Emit(", ")
		}
		parameterTypeName := updateTypeName(parameter.Type)
		generator.Emitf("%s %s", parameterTypeName.FullName(), parameter.Name.Value)
	}
	generator.Emit(")")
	return generation.buffer.String()
}

func (generation *MethodGeneration) generateBody() string {
	generation.buffer.Reset()
	// Do not use the normal BlockStatement visitor for BlockStatement generation in
	// a method body. It will generate open and close brackets and this will produce
	// faulty code, if a prologue or epilogue is generated.
	if block, ok := generation.declaration.Body.(*ast.BlockStatement); ok {
		for _, child := range block.Children {
			generation.generator.EmitNode(child)
		}
		return generation.buffer.String()
	}
	generation.generator.EmitNode(generation.declaration.Body)
	return generation.buffer.String()
}

func (generation *MethodGeneration) generatePrologue() string {
	generation.buffer.Reset()
	for _, prologueGenerator := range generation.prologueGenerators {
		prologueGenerator()
	}
	if len(generation.prologueGenerators) > 0 {
		generation.generator.Emit("\n")
		generation.generator.Spaces()
	}
	return generation.buffer.String()
}

func (generation *MethodGeneration) generateEpilogue() string {
	generation.buffer.Reset()
	for _, epilogueGenerator := range generation.epilogueGenerators {
		epilogueGenerator()
	}
	return generation.buffer.String()
}

func (generation *MethodGeneration) addEpilogueGenerator(name string, function EpilogueGenerator) {
	generation.epilogueGenerators[name] = function
}

func (generation *MethodGeneration) addPrologueGenerator(name string, function PrologueGenerator) {
	generation.prologueGenerators[name] = function
}
