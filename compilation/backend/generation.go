package backend

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"strings"
)

// Generation generates C code from an ast.
type Generation struct {
	unit                        *ast.TranslationUnit
	output                      *strings.Builder
	buffer                      *strings.Builder
	method                      *MethodDefinition
	visitor                     *ast.Visitor
	indent                      int8
	appendNewLineAfterStatement bool
	importModules               map[string]string
}

type FileNaming interface {
	FileNameForUnit(unit *ast.TranslationUnit)
}

func NewGenerationWithExtension(unit *ast.TranslationUnit, extension Extension) *Generation {
	generation := NewGeneration(unit)
	extension.ModifyVisitor(generation, generation.visitor)
	return generation
}

// NewGeneration constructs a Generation that generates C code from
// the nodes in the passed translation-unit.
func NewGeneration(unit *ast.TranslationUnit) (generation *Generation) {
	generation = &Generation{
		unit:                        unit,
		output:                      &strings.Builder{},
		importModules:               map[string]string{},
		appendNewLineAfterStatement: true,
	}
	generation.buffer = generation.output
	generation.visitor = CreateGenericCppVisitor(generation)
	return
}

func (generation *Generation) Filename() string {
	return fmt.Sprintf("%s.cc", generation.unit.ToTypeName().NonGenericName())
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

func (generation *Generation) EmitNode(node ast.Node) {
	node.Accept(generation.visitor)
}

func (generation *Generation) Generate() string {
	generation.EmitNode(generation.unit)
	return generation.String()
}
