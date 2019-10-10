package tree

import "gitlab.com/strict-lang/sdk/pkg/compilation/input"

// Node is implemented by every node of the tree.
type Node interface {
	Locate() input.Region
	// Accept makes the visitor visit this node.
	Accept(visitor Visitor)
	// AcceptRecursive makes the visitor visit this node and its children.
	AcceptRecursive(visitor Visitor)
}

// Named is implemented by all nodes that have a name.
type Named interface {
	// Name returns the nodes name.
	Name() string
}

type NodeKind int

const (
	invalidKind NodeKind = iota
	expressionKindBegin
	IdentifierNodeKind
	StringLiteralNodeKind
	NumberLiteralNodeKind
	FieldSelectExpressionNodeKind
	ListSelectExpressionNodeKind
	BinaryExpressionNodeKind
	UnaryExpressionNodeKind
	PostfixExpressionNodeKind
	CreateExpressionNodeKind
	CallArgumentNodeKind
	CallExpressionNodeKind
	expressionKindEnd
	statementKindBegin
	ConditionalStatementNodeKind
	InvalidStatementNodeKind
	YieldStatementNodeKind
	BlockStatementNodeKind
	AssertStatementNodeKind
	ReturnStatementNodeKind
	ImportStatementNodeKind
	EmptyStatementNodeKind
	TestStatementNodeKind
	AssignStatementNodeKind
	ExpressionStatementNodeKind
	ForEachLoopStatementNodeKind
	RangedLoopStatementNodeKind
	statementKindEnd
	declarationKindBegin
	ParameterNodeKind
	FieldDeclarationNodeKind
	MethodDeclarationNodeKind
	ClassDeclarationNodeKind
	ConstructorDeclarationNodeKind
	declarationKindEnd
	typeNameKindBegin
	ListTypeNameNodeKind
	GenericTypeNameNodeKind
	ConcreteTypeNameNodeKind
	typeNameKindEnd
	TranslationUnitNodeKind
	WildcardNodeKind
)

// IsExpression returns true if the kind is an expression.
func (kind NodeKind) IsExpression() bool {
	return kind.isInExclusiveRange(kind, expressionKindBegin, expressionKindEnd)
}

// IsStatement returns true if the kind is a statement.
func (kind NodeKind) IsStatement() bool {
	return kind.isInExclusiveRange(kind, statementKindBegin, statementKindEnd)
}

// IsDeclaration returns true if the kind is a declaration.
func (kind NodeKind) IsDeclaration() bool {
	return kind.isInExclusiveRange(kind, declarationKindBegin, declarationKindEnd)
}

func (kind NodeKind) isInExclusiveRange(tested, begin, end NodeKind) bool {
	return tested > begin && tested < end
}
