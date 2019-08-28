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
	method                      *MethodGeneration
	visitor                     *ast.Visitor
	indent                      int8
	appendNewLineAfterStatement bool
	importModules               map[string]string
	settings                    Settings
}

type FileNaming interface {
	FileNameForUnit(unit *ast.TranslationUnit)
}

// NewCodeGenerator constructs a Generation that generates C code from
// the nodes in the passed translation-unit.
func NewCodeGenerator(settings Settings, unit *ast.TranslationUnit) (generator *Generation) {
	generator = &Generation{
		unit:                        unit,
		output:                      &strings.Builder{},
		settings:                    settings,
		importModules:               map[string]string{},
		appendNewLineAfterStatement: true,
	}
	generator.buffer = generator.output
	generator.visitor = CreateGenericCppVisitor(generator)
	return
}

func (generation *Generation) String() string {
	return generation.output.String()
}

func (generation *Generation) Emit(code string) {
	generation.buffer.WriteString(code)
}

func (generation *Generation) Emitf(code string, arguments ...interface{}) {
	formatted := fmt.Sprintf(code, arguments...)
	generation.buffer.WriteString(formatted)
}

func (generation *Generation) enterBlock() {
	generation.indent++
}

func (generation *Generation) leaveBlock() {
	generation.indent--
	if generation.indent < 0 {
		generation.indent = 0
	}
}

func (generation *Generation) writeEndOfStatement() {
	if generation.appendNewLineAfterStatement {
		generation.Emit("\n")
	}
}

func (generation *Generation) Spaces() {
	for index := int8(0); index < generation.indent; index++ {
		generation.Emit("  ")
	}
}

func (generation *Generation) EmitNode(node ast.Node) {
	node.Accept(generation.visitor)
}

func (generation *Generation) Generate() string {
	generation.EmitNode(generation.unit)
	return generation.String()
}
