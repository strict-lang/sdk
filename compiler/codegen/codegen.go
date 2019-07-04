package codegen

import (
	"fmt"
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"strings"
)

// CodeGenerator generates C code from an ast.
type CodeGenerator struct {
	unit       *ast.TranslationUnit
	output     *strings.Builder
	buffer     *strings.Builder
	method     *MethodGeneration
	generators *ast.Visitor
	indent     int8
}

// NewCodeGenerator constructs a CodeGenerator that generates C code from
// the nodes in the passed translation-unit.
func NewCodeGenerator(unit *ast.TranslationUnit) *CodeGenerator {
	generators := ast.NewEmptyVisitor()
	codeGenerator := &CodeGenerator{
		unit:       unit,
		output:     &strings.Builder{},
		generators: generators,
	}
	codeGenerator.buffer = codeGenerator.output
	generators.VisitMethod = codeGenerator.GenerateMethod
	generators.VisitIdentifier = codeGenerator.GenerateIdentifier
	generators.VisitMethodCall = codeGenerator.GenerateMethodCall
	generators.VisitStringLiteral = codeGenerator.GenerateStringLiteral
	generators.VisitNumberLiteral = codeGenerator.GenerateNumberLiteral
	generators.VisitYieldStatement = codeGenerator.GenerateYieldStatement
	generators.VisitBlockStatement = codeGenerator.GenerateBlockStatement
	generators.VisitReturnStatement = codeGenerator.GenerateReturnStatement
	generators.VisitTranslationUnit = codeGenerator.GenerateTranslationUnit
	generators.VisitUnaryExpression = codeGenerator.GenerateUnaryExpression
	generators.VisitBinaryExpression = codeGenerator.GenerateBinaryExpression
	generators.VisitExpressionStatement = codeGenerator.GenerateExpressionStatement
	generators.VisitFromToLoopStatement = codeGenerator.GenerateFromToLoopStatement
	generators.VisitConditionalStatement = codeGenerator.GenerateConditionalStatement
	generators.VisitForeachLoopStatement = codeGenerator.GenerateForEachLoopStatement
	return codeGenerator
}

func (generator *CodeGenerator) String() string {
	return generator.output.String()
}

func (generator *CodeGenerator) Emit(code string) {
	generator.buffer.WriteString(code)
}

func (generator *CodeGenerator) Emitf(code string, arguments ...interface{}) {
	formatted := fmt.Sprintf(code, arguments...)
	generator.buffer.WriteString(formatted)
}

func (generator *CodeGenerator) enterBlock() {
	generator.indent++
}

func (generator *CodeGenerator) leaveBlock() {
	generator.indent--
	if generator.indent < 0 {
		generator.indent = 0
	}
}

func (generator *CodeGenerator) Spaces() {
	for index := int8(0); index < generator.indent; index++ {
		generator.Emit("  ")
	}
}

func (generator *CodeGenerator) Generate() string {
	generator.unit.Accept(generator.generators)
	return generator.String()
}
