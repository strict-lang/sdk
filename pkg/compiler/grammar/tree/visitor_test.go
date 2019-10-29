package tree

import (
	"testing"
)

type VisitorTest struct {
	visitor       Visitor
	expectedKinds []NodeKind
	testing       *testing.T
	tested        Node
}

func CreateVisitorTest(node Node, testing *testing.T) *VisitorTest {
	test := &VisitorTest{
		tested:  node,
		testing: testing,
	}
	test.visitor = test.createVisitor()
	return test
}

func (test *VisitorTest) Run() {
	test.tested.Accept(test.visitor)
	test.ensureComplete()
}

func (test *VisitorTest) RunRecursive() {
	test.tested.AcceptRecursive(test.visitor)
	test.ensureComplete()
}

func (test *VisitorTest) Expect(kind NodeKind) *VisitorTest {
	test.expectedKinds = append(test.expectedKinds, kind)
	return test
}

func (test *VisitorTest) ensureComplete() {
	if len(test.expectedKinds) != 0 {
		test.testing.Error("Not every expected NodeKind has been visited")
		for _, remaining := range test.expectedKinds {
			test.testing.Errorf("- %s has not been visisted", remaining)
		}
	}
}

func (test *VisitorTest) popExpectedNodeKind() (NodeKind, bool) {
	expectedKindsCount := len(test.expectedKinds)
	if expectedKindsCount == 0 {
		return invalidKind, false
	}
	next := test.expectedKinds[0]
	test.expectedKinds = test.expectedKinds[1:expectedKindsCount]
	return next, true
}

func (test *VisitorTest) reportNodeEncounter(kind NodeKind) {
	expected, isExpectingNode := test.popExpectedNodeKind()
	if !isExpectingNode {
		test.testing.Error("Visited more nodes than expected")
		return
	}
	if expected != kind {
		test.testing.Errorf(`Visited unexpected node.
  Expected: %s
  Received: %s
`, expected, kind)
	}
}

func (test *VisitorTest) createVisitor() Visitor {
	return &DelegatingVisitor{
		ParameterVisitor: func(*Parameter) {
			test.reportNodeEncounter(ParameterNodeKind)
		},
		IdentifierVisitor: func(*Identifier) {
			test.reportNodeEncounter(IdentifierNodeKind)
		},
		CallArgumentVisitor: func(*CallArgument) {
			test.reportNodeEncounter(CallArgumentNodeKind)
		},
		ListTypeNameVisitor: func(*ListTypeName) {
			test.reportNodeEncounter(ListTypeNameNodeKind)
		},
		TestStatementVisitor: func(*TestStatement) {
			test.reportNodeEncounter(TestStatementNodeKind)
		},
		StringLiteralVisitor: func(*StringLiteral) {
			test.reportNodeEncounter(StringLiteralNodeKind)
		},
		NumberLiteralVisitor: func(*NumberLiteral) {
			test.reportNodeEncounter(NumberLiteralNodeKind)
		},
		CallExpressionVisitor: func(*CallExpression) {
			test.reportNodeEncounter(CallExpressionNodeKind)
		},
		EmptyStatementVisitor: func(*EmptyStatement) {
			test.reportNodeEncounter(EmptyStatementNodeKind)
		},
		YieldStatementVisitor: func(*YieldStatement) {
			test.reportNodeEncounter(YieldStatementNodeKind)
		},
		BlockStatementVisitor: func(*BlockStatement) {
			test.reportNodeEncounter(BlockStatementNodeKind)
		},
		AssertStatementVisitor: func(*AssertStatement) {
			test.reportNodeEncounter(AssertStatementNodeKind)
		},
		UnaryExpressionVisitor: func(*UnaryExpression) {
			test.reportNodeEncounter(UnaryExpressionNodeKind)
		},
		ImportStatementVisitor: func(*ImportStatement) {
			test.reportNodeEncounter(ImportStatementNodeKind)
		},
		AssignStatementVisitor: func(*AssignStatement) {
			test.reportNodeEncounter(AssignStatementNodeKind)
		},
		ReturnStatementVisitor: func(*ReturnStatement) {
			test.reportNodeEncounter(ReturnStatementNodeKind)
		},
		TranslationUnitVisitor: func(*TranslationUnit) {
			test.reportNodeEncounter(TranslationUnitNodeKind)
		},
		CreateExpressionVisitor: func(*CreateExpression) {
			test.reportNodeEncounter(CreateExpressionNodeKind)
		},
		InvalidStatementVisitor: func(*InvalidStatement) {
			test.reportNodeEncounter(InvalidStatementNodeKind)
		},
		FieldDeclarationVisitor: func(*FieldDeclaration) {
			test.reportNodeEncounter(FieldDeclarationNodeKind)
		},
		GenericTypeNameVisitor: func(*GenericTypeName) {
			test.reportNodeEncounter(GenericTypeNameNodeKind)
		},
		ConcreteTypeNameVisitor: func(*ConcreteTypeName) {
			test.reportNodeEncounter(ConcreteTypeNameNodeKind)
		},
		ClassDeclarationVisitor: func(*ClassDeclaration) {
			test.reportNodeEncounter(ClassDeclarationNodeKind)
		},
		BinaryExpressionVisitor: func(*BinaryExpression) {
			test.reportNodeEncounter(BinaryExpressionNodeKind)
		},
		MethodDeclarationVisitor: func(*MethodDeclaration) {
			test.reportNodeEncounter(MethodDeclarationNodeKind)
		},
		RangedLoopStatementVisitor: func(*RangedLoopStatement) {
			test.reportNodeEncounter(RangedLoopStatementNodeKind)
		},
		ExpressionStatementVisitor: func(*ExpressionStatement) {
			test.reportNodeEncounter(ExpressionStatementNodeKind)
		},
		ForEachLoopStatementVisitor: func(*ForEachLoopStatement) {
			test.reportNodeEncounter(ForEachLoopStatementNodeKind)
		},
		ConditionalStatementVisitor: func(*ConditionalStatement) {
			test.reportNodeEncounter(ConditionalStatementNodeKind)
		},
		ListSelectExpressionVisitor: func(*ListSelectExpression) {
			test.reportNodeEncounter(ListSelectExpressionNodeKind)
		},
		ConstructorDeclarationVisitor: func(*ConstructorDeclaration) {
			test.reportNodeEncounter(ConstructorDeclarationNodeKind)
		},
		FieldSelectExpressionVisitor: func(*FieldSelectExpression) {
			test.reportNodeEncounter(FieldSelectExpressionNodeKind)
		},
		PostfixExpressionVisitor: func(*PostfixExpression) {
			test.reportNodeEncounter(PostfixExpressionNodeKind)
		},
		WildcardNodeVisitor: func(*WildcardNode) {
			test.reportNodeEncounter(WildcardNodeKind)
		},
	}
}
