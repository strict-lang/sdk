package tree

import "gitlab.com/strict-lang/sdk/pkg/compiler/input"

type TestStatement struct {
	MethodName string
	Child      Node
	Region     input.Region
}

func (test *TestStatement) Accept(visitor Visitor) {
	VisitTestStatement(test)
}

func (test *TestStatement) AcceptRecursive(visitor Visitor) {
	VisitTestStatement(test)
	AcceptRecursive(visitor)
}

func (test *TestStatement) Locate() input.Region {
	return test.Region
}