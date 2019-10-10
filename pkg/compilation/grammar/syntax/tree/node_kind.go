package tree

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