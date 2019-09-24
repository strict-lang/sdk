package testfile

import (
	"fmt"
	backend2 "gitlab.com/strict-lang/sdk/pkg/compilation/backend"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

type Extension struct{}

var _ backend2.Extension = &Extension{}

type TestFile struct {
	generation *backend2.Generation
}

func (extension *Extension) ModifyVisitor(generation *backend2.Generation, visitor *syntaxtree2.Visitor) {
	testFile := &TestFile{
		generation: generation,
	}
	visitor.VisitAssertStatement = testFile.emitAssertStatement
}

const typeTestingInstance = "testing"
const methodTestingInstance = "methodTesting"
const returnOnFailure = false

func (testFile *TestFile) emitAssertStatement(statement *syntaxtree2.AssertStatement) {
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

func (testFile *TestFile) emitFailedAssertion(statement *syntaxtree2.AssertStatement) {
	failureMessage := generateAssertionFailureMessage(statement.Expression)
	testFile.generation.EmitFormatted("%s.ReportFailedAssertion(\"%s\");",
		methodTestingInstance, failureMessage)
	if returnOnFailure {
		testFile.generation.Emit("return;")
	}
}

func generateAssertionFailureMessage(expression syntaxtree2.Node) string {
	assertionMessage := backend2.NewAssertionMessageComputation()
	assertionMessage.GenerateNode(expression)
	return assertionMessage.String()
}

type TestDefinition struct {
	testedMethodName string
	testMethodName   string
	node             syntaxtree2.TestStatement
	generation       *backend2.Generation
}

func NewTestDefinition(node syntaxtree2.TestStatement, generation *backend2.Generation) *TestDefinition {
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
