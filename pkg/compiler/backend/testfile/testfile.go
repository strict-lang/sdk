package testfile

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/backend"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Extension struct{}

var _ backend.Extension = &Extension{}

type TestFile struct {
	generation *backend.Generation
}

func NewGeneration() *Extension {
	return &Extension{}
}

func (extension *Extension) ModifyVisitor(generation *backend.Generation, visitor *tree.DelegatingVisitor) {
	testFile := &TestFile{
		generation: generation,
	}
	visitor.AssertStatementVisitor = testFile.emitAssertStatement
	visitor.TestStatementVisitor = testFile.emitTestDefinition
	visitor.MethodDeclarationVisitor = testFile.emitMethodDeclaration
}

const typeTestingInstance = "testing"
const methodTestingInstance = "methodTesting"
const returnOnFailure = false

func (testFile *TestFile) emitAssertStatement(statement *tree.AssertStatement) {
	generation := testFile.generation
	generation.Emit("if (!(")
	generation.EmitNode(statement.Expression)
	generation.Emit(")) {\n")
	generation.IncreaseIndent()
	generation.EmitIndent()
	testFile.emitFailedAssertion(statement)
	generation.Emit("\n")
	generation.DecreaseIndent()
	generation.EmitIndent()
	generation.Emit("}")
}

func (testFile *TestFile) emitTestDefinition(node *tree.TestStatement) {
	test := NewTestDefinition(node, testFile.generation)
	test.Emit()
}

func (testFile *TestFile) emitFailedAssertion(statement *tree.AssertStatement) {
	failureMessage := generateAssertionFailureMessage(statement.Expression)
	testFile.generation.EmitFormatted("%s.ReportFailedAssertion(\"%s\");",
		methodTestingInstance, failureMessage)
	if returnOnFailure {
		testFile.generation.Emit("return;")
	}
}

func (testFile *TestFile) emitMethodDeclaration(method *tree.MethodDeclaration) {
	block, isBlock := method.Body.(*tree.StatementBlock)
	if !isBlock {
		return
	}
	for _, statement := range block.Children {
		if test, isTest := statement.(*tree.TestStatement); isTest {
			testFile.generation.EmitNode(test)
		}
	}
}

func generateAssertionFailureMessage(expression tree.Node) string {
	assertionMessage := backend.NewAssertionMessageComputation()
	assertionMessage.GenerateNode(expression)
	return assertionMessage.String()
}

type TestDefinition struct {
	testedMethodName string
	testMethodName   string
	node             *tree.TestStatement
	generation       *backend.Generation
}

func NewTestDefinition(node *tree.TestStatement, generation *backend.Generation) *TestDefinition {
	return &TestDefinition{
		testedMethodName: node.MethodName,
		testMethodName:   produceTestMethodName(node.MethodName),
		node:             node,
		generation:       generation,
	}
}

func (node *TestDefinition) Emit() {
	generation := node.generation
	generation.EmitFormatted(
		"void %s(testing::Testing* %s) {",
		node.testMethodName,
		typeTestingInstance)
	generation.IncreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
	node.emitMethodTestingRAII()
	generation.EmitIndent()
	generation.EmitNode(node.node.Body)
	generation.DecreaseIndent()
	generation.EmitEndOfLine()
	generation.EmitIndent()
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
