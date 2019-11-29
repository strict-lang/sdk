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
	test.visitor = NewReportingVisitor(test)
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
