package testfile

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/backend"
)

type Extension struct{}

var _ backend.Extension = &Extension{}

type TestFile struct {
	generation *backend.Generation
}

func (extension *Extension) ModifyVisitor(generation *backend.Generation, visitor *ast.Visitor) {
	testFile := &TestFile{
		generation: generation,
	}
	visitor.VisitAssertStatement = testFile.emitAssertStatement
}

const typeTestingInstance = "testing"
const methodTestingInstance = "methodTesting"
const returnOnFailure = false

func (testFile *TestFile) emitAssertStatement(statement *ast.AssertStatement) {
	generation := testFile.generation
	generation.Emit("if (!(")
	generation.EmitNode(statement.Expression)
	generation.Emit(") {")
	generation.IncreaseIndent()
	generation.EmitIndent()
	testFile.emitFailedAssertion(statement)
	generation.EmitEndOfLine()
	generation.DecreaseIndent()
	generation.Emit("}")
	generation.EmitEndOfLine()
}

func (testFile *TestFile) emitFailedAssertion(statement *ast.AssertStatement) {
	failureMessage := generateAssertionFailureMessage(statement.Expression)
	testFile.generation.EmitFormatted("%s.ReportFailedAssertion(\"%s\");",
		methodTestingInstance, failureMessage)
	if returnOnFailure {
		testFile.generation.Emit("return;")
	}
}

func generateAssertionFailureMessage(expression ast.Node) string {
	assertionMessage := backend.newAssertionMessageComputation()
	assertionMessage.generateNode(expression)
	return assertionMessage.String()
}

type TestDefinition struct {
	testedMethodName string
	testMethodName   string
	node             ast.TestStatement
	generation       *backend.Generation
}

func NewTestDefinition(node ast.TestStatement, generation *backend.Generation) *TestDefinition {
	return &TestDefinition{
		testedMethodName: node.MethodName,
		testMethodName:   produceTestMethodName(node.MethodName),
		node:             node,
		generation:       generation,
	}
}

func (node *TestDefinition) Emit() {
	generation := node.generation
	generation.EmitFormatted("void %s(%s *testing::Testing) {", typeTestingInstance, node.testMethodName)
	generation.IncreaseIndent()
	generation.EmitIndent()
	node.emitMethodTestingRAII()
	generation.DecreaseIndent()
	generation.Emit("}")
}

func (node *TestDefinition) emitMethodTestingRAII() {
	node.generation.EmitFormatted("MethodTesting %s(%s, \"%s\");",
		methodTestingInstance, typeTestingInstance, node.testedMethodName)
	node.generation.EmitEndOfLine()
}

func produceTestMethodName(methodName string) string {
	return fmt.Sprintf("Test_%s_", methodName)
}
