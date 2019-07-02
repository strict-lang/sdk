package codegen

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
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
	generator.leaveBlock()

	generator.buffer = generator.output
	generator.Emit(declaration)
	generator.Emit(" {\n")
	generator.Emit(prologue)
	generator.Emit(methodBody)
	generator.Emit(epilogue)
	generator.Emit("\n")
	generator.Spaces()
	generator.Emit("}\n\n")
}

func (generation *MethodGeneration) generateDeclaration() string {
	generation.buffer.Reset()
	generator := generation.generator

	method := generation.declaration
	returnTypeName := updateTypeName(method.Type)

	generator.Spaces()
	generator.Emitf("%s %s(", returnTypeName.FullName(), method.Name.Value)
	for _, parameter := range method.Parameters {
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
			child.Accept(generation.generator.generators)
		}
		return generation.buffer.String()
	}
	generation.declaration.Body.Accept(generation.generator.generators)
	return generation.buffer.String()
}

func (generation *MethodGeneration) generatePrologue() string {
	generation.buffer.Reset()
	for _, prologueGenerator := range generation.prologueGenerators {
		prologueGenerator()
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
