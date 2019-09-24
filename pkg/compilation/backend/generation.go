package backend

import (
	"fmt"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	"strings"
)

// Generation generates C code from an syntaxtree.
type Generation struct {
	Unit                          *syntaxtree.TranslationUnit
	output                        *strings.Builder
	buffer                        *strings.Builder
	method                        *MethodDefinition
	visitor                       *syntaxtree.Visitor
	indent                        int8
	appendNewLineAfterStatement   bool
	importModules                 map[string]string
	shouldInsertNamespaceSelector bool
}

type FileNaming interface {
	FileNameForUnit(unit *syntaxtree.TranslationUnit)
}

func NewGenerationWithExtension(unit *syntaxtree.TranslationUnit, extension Extension) *Generation {
	generation := NewGeneration(unit)
	extension.ModifyVisitor(generation, generation.visitor)
	return generation
}

// NewGeneration constructs a Generation that generates C code from
// the nodes in the passed translation-Unit.
func NewGeneration(unit *syntaxtree.TranslationUnit) (generation *Generation) {
	generation = &Generation{
		Unit:                          unit,
		output:                        &strings.Builder{},
		importModules:                 map[string]string{},
		appendNewLineAfterStatement:   true,
		shouldInsertNamespaceSelector: true,
	}
	generation.buffer = generation.output
	generation.visitor = CreateGenericCppVisitor(generation)
	return
}

func (generation *Generation) DisableNamespaceSelectors() {
	generation.shouldInsertNamespaceSelector = false
}

func (generation *Generation) Filename() string {
	return fmt.Sprintf("%s.cc", generation.Unit.ToTypeName().NonGenericName())
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

func (generation *Generation) EmitNode(node syntaxtree.Node) {
	node.Accept(generation.visitor)
}

func (generation *Generation) Generate() string {
	generation.EmitNode(generation.Unit)
	return generation.String()
}
