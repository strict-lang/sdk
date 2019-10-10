package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

type TestStatement struct {
	MethodName   string
	Child   Node
	Region input.Region
}

func (test *TestStatement) Accept(visitor Visitor) {
	visitor.VisitTestStatement(test)
}

func (test *TestStatement) AcceptRecursive(visitor Visitor) {
	visitor.VisitTestStatement(test)
	test.Child.AcceptRecursive(visitor)
}

func (test *TestStatement) Locate() input.Region {
	return test.Region
}