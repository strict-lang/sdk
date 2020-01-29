package cpp

import (
	"fmt"
	"strict.dev/sdk/pkg/compiler/backend"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strings"
)

// HeaderFileGeneration generates C code from an tree.
type Generation struct {
	Unit                          *tree.TranslationUnit
	output                        *strings.Builder
	buffer                        *strings.Builder
	method                        *MethodDefinition
	visitor                       *tree.DelegatingVisitor
	indent                        int8
	appendNewLineAfterStatement   bool
	importModules                 map[string]string
	shouldInsertNamespaceSelector bool
	shouldImportStdlibClasses     bool
}

type FileNaming interface {
	FileNameForUnit(unit *tree.TranslationUnit)
}

func NewGenerationWithExtension(input backend.Input, extension Extension) *Generation {
	generation := NewGeneration(input.Unit)
	extension.ModifyVisitor(generation, generation.visitor)
	return generation
}

// NewGeneration constructs a HeaderFileGeneration that generates C code from
// the nodes in the passed translation-Unit.
func NewGeneration(unit *tree.TranslationUnit) (generation *Generation) {
	generation = &Generation{
		Unit:                          unit,
		output:                        &strings.Builder{},
		importModules:                 map[string]string{},
		appendNewLineAfterStatement:   true,
		shouldInsertNamespaceSelector: true,
		shouldImportStdlibClasses:     true,
	}
	generation.buffer = generation.output
	generation.visitor = CreateGenericCppVisitor(generation)
	return
}

func (generation *Generation) DisableStdlibClassImport() {
	generation.shouldImportStdlibClasses = false
}

func (generation *Generation) DisableNamespaceSelectors() {
	generation.shouldInsertNamespaceSelector = false
}

func (generation *Generation) Filename() string {
	return fmt.Sprintf("%s.cc", generation.Unit.ToTypeName().BaseName())
}

func (generation *Generation) String() string {
	return generation.output.String()
}

func (generation *Generation) Emit(code string) {
	generation.buffer.WriteString(code)
}

func (generation *Generation) EmitIndent() {
	for index := int8(0); index < generation.indent; index++ {
		generation.Emit(" ")
	}
}

func (generation *Generation) EmitFormatted(code string, arguments ...interface{}) {
	formatted := fmt.Sprintf(code, arguments...)
	generation.buffer.WriteString(formatted)
}

func (generation *Generation) EmitEndOfLine() {
	if generation.appendNewLineAfterStatement {
		generation.Emit("\n")
	}
}

func (generation *Generation) IncreaseIndent() {
	generation.indent++
}

func (generation *Generation) DecreaseIndent() {
	generation.indent--
	if generation.indent < 0 {
		generation.indent = 0
	}
}

func (generation *Generation) EmitNode(node tree.Node) {
	node.Accept(generation.visitor)
}

func (generation *Generation) Generate() string {
	generation.EmitNode(generation.Unit)
	return generation.String()
}
