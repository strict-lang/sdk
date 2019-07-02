package codegen

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/pkg/ast"
	"strings"
)

// CodeGenerator generates C code from an ast.
type CodeGenerator struct {
	unit   		 *ast.TranslationUnit
	output 		 *strings.Builder
	buffer     *strings.Builder
	method     *MethodGeneration
	generators *ast.Visitor
}

// NewCodeGenerator constructs a CodeGenerator that generates C code from
// the nodes in the passed translation-unit.
func NewCodeGenerator(unit *ast.TranslationUnit) *CodeGenerator {
	generators := ast.NewEmptyVisitor()
	codeGenerator := &CodeGenerator{
		unit:   unit,
		output: &strings.Builder{},
		generators: generators,
	}
	codeGenerator.buffer = codeGenerator.output
	generators.VisitMethod = codeGenerator.GenerateMethod
	generators.VisitIdentifier = codeGenerator.GenerateIdentifier
	generators.VisitMethodCall = codeGenerator.GenerateMethodCall
	generators.VisitStringLiteral = codeGenerator.GenerateStringLiteral
	generators.VisitNumberLiteral = codeGenerator.GenerateNumberLiteral
	generators.VisitYieldStatement = codeGenerator.GenerateYieldStatement
	generators.VisitConditionalStatement = codeGenerator.GenerateConditionalStatement
	return codeGenerator
}

func (generator *CodeGenerator) String() string {
	return generator.output.String()
}

func (generator *CodeGenerator) Emit(code string) {
	generator.output.WriteString(code)
}

func (generator *CodeGenerator) Emitf(code string, arguments ...interface{}) {
	formatted := fmt.Sprintf(code, arguments)
	generator.output.WriteString(formatted)
}
