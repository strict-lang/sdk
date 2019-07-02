package codegen

import (
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"strings"
)

type MethodGeneration struct {
	generator *CodeGenerator
	prologueGenerators []PrologueGenerator
	epilogueGenerators []EpilogueGenerator
	declaration ast.Method
	buffer *strings.Builder
}

type PrologueGenerator func()
type EpilogueGenerator func()

func (generator *CodeGenerator) NewMethodGeneration(method ast.Method) *MethodGeneration {
	return &MethodGeneration{
		generator: generator,
		declaration: method,
		buffer: &strings.Builder{},
		epilogueGenerators: []EpilogueGenerator{},
		prologueGenerators: []PrologueGenerator{},
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
	methodBody := generation.generateBody()
	prologue := generation.generatePrologue()
	epilogue := generation.generateEpilogue()

	generation.buffer = generator.output
	generation.buffer.WriteString(declaration)
	generation.buffer.WriteString("\n{")
	generation.buffer.WriteString(prologue)
	generation.buffer.WriteString(methodBody)
	generation.buffer.WriteString(epilogue)
	generation.buffer.WriteString("}")
}

func (generation *MethodGeneration) generateDeclaration() string {
	generation.buffer.Reset()

	method := generation.declaration
	updateTypeName(method.Type)
	name := method.Type.FullName()

	generation.buffer.WriteString(name)
	generation.buffer.WriteString(" ")
	generation.buffer.WriteString(method.Name.Value)
	generation.buffer.WriteString("(")
	for _, parameter := range method.Parameters {
		updateTypeName(parameter.Type)
		generation.buffer.WriteString(parameter.Type.FullName())
		generation.buffer.WriteString( " ")
		generation.buffer.WriteString(parameter.Name.Value)
	}
	generation.buffer.WriteString(")")
	return generation.buffer.String()
}

func (generation *MethodGeneration) generateBody() string {
	generation.buffer.Reset()
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

func (generation *MethodGeneration) addEpilogueGenerator(function EpilogueGenerator) {
	generation.epilogueGenerators = append(generation.epilogueGenerators, function)
}

func (generation *MethodGeneration) addPrologueGenerator(function PrologueGenerator) {
	generation.prologueGenerators = append(generation.prologueGenerators, function)
}
