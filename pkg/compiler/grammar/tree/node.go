package tree

import "github.com/strict-lang/sdk/pkg/compiler/input"

// Node is implemented by every node of the tree.
type Node interface {
	Locate() input.Region
	// Accept makes the visitor visit this node.
	Accept(visitor Visitor)
	// AcceptRecursive makes the visitor visit this node and its children.
	AcceptRecursive(visitor Visitor)
	// Matches checks whether the instance matches the passed node. It does not
	// take positions into account.
	Matches(node Node) bool
	// Enclosing returns the node that encloses the passed node. In the case
	// of parameters this is the method they belong to.
	EnclosingNode() (node Node, exists bool)
	SetEnclosingNode(target Node)
}

// Named is implemented by all nodes that have a name.
type Named interface {
	// Name returns the nodes name.
	Name() string
}

type NodeKind int

const (
	invalidKind NodeKind = iota
	UnknownNodeKind
	expressionKindBegin
	IdentifierNodeKind
	StringLiteralNodeKind
	NumberLiteralNodeKind
	ChainExpressionNodeKind
	ListSelectExpressionNodeKind
	BinaryExpressionNodeKind
	UnaryExpressionNodeKind
	PostfixExpressionNodeKind
	CreateExpressionNodeKind
	CallArgumentNodeKind
	CallExpressionNodeKind
	LetBindingNodeKind
	expressionKindEnd
	statementKindBegin
	ConditionalStatementNodeKind
	InvalidStatementNodeKind
	BreakStatementNodeKind
	YieldStatementNodeKind
	StatementBlockNodeKind
	AssertStatementNodeKind
	ReturnStatementNodeKind
	ImportStatementNodeKind
	EmptyStatementNodeKind
	TestStatementNodeKind
	AssignStatementNodeKind
	ExpressionStatementNodeKind
	ForEachLoopStatementNodeKind
	RangedLoopStatementNodeKind
	ImplementStatementNodeKind
	ListExpressionNodeKind
	GenericStatementNodeKind
	statementKindEnd
	declarationKindBegin
	ParameterNodeKind
	FieldDeclarationNodeKind
	MethodDeclarationNodeKind
	ClassDeclarationNodeKind
	ConstructorDeclarationNodeKind
	declarationKindEnd
	typeNameKindBegin
	TypeNameNodeGroup // Used only in parsing
	ListTypeNameNodeKind
	GenericTypeNameNodeKind
	ConcreteTypeNameNodeKind
	OptionalTypeNameNodeKind
	typeNameKindEnd
	TranslationUnitNodeKind
	WildcardNodeKind
)

var nodeKindNames = map[NodeKind]string{
	UnknownNodeKind:                "Unknown",
	IdentifierNodeKind:             "Identifier",
	StringLiteralNodeKind:          "StringLiteral",
	NumberLiteralNodeKind:          "NumberLiteral",
	ChainExpressionNodeKind:        "ChainExpression",
	ListSelectExpressionNodeKind:   "ListSelectExpression",
	BinaryExpressionNodeKind:       "BinaryExpression",
	UnaryExpressionNodeKind:        "UnaryExpression",
	PostfixExpressionNodeKind:      "PostfixExpression",
	CreateExpressionNodeKind:       "CreateExpression",
	CallArgumentNodeKind:           "CallArgument",
	CallExpressionNodeKind:         "CallExpression",
	ConditionalStatementNodeKind:   "ConditionalStatement",
	InvalidStatementNodeKind:       "InvalidStatement",
	YieldStatementNodeKind:         "YieldStatement",
	StatementBlockNodeKind:         "StatementBlock",
	AssertStatementNodeKind:        "AssertStatement",
	ReturnStatementNodeKind:        "ReturnStatement",
	ImportStatementNodeKind:        "ImportStatement",
	EmptyStatementNodeKind:         "EmptyStatement",
	BreakStatementNodeKind:         "BreakStatement",
	TestStatementNodeKind:          "TestStatement",
	AssignStatementNodeKind:        "AssignStatement",
	ExpressionStatementNodeKind:    "ExpressionStatement",
	ForEachLoopStatementNodeKind:   "ForEachLoopStatement",
	RangedLoopStatementNodeKind:    "RangedLoopStatement",
	ParameterNodeKind:              "Parameter",
	FieldDeclarationNodeKind:       "FieldDeclaration",
	MethodDeclarationNodeKind:      "MethodDeclaration",
	ClassDeclarationNodeKind:       "ClassDeclaration",
	ConstructorDeclarationNodeKind: "ConstructorDeclaration",
	TypeNameNodeGroup:              "TypeName",
	ListTypeNameNodeKind:           "ListTypeName",
	GenericTypeNameNodeKind:        "GenericTypeName",
	ConcreteTypeNameNodeKind:       "ConcreteTypeName",
	OptionalTypeNameNodeKind:       "OptionalTypeName",
	TranslationUnitNodeKind:        "TranslationUnit",
	LetBindingNodeKind:             "LetBinding",
	ImplementStatementNodeKind:     "ImplementStatement",
	GenericStatementNodeKind:       "GenericStatement",
	WildcardNodeKind:               "Wildcard",
}

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

func (kind NodeKind) Name() string {
	name, ok := nodeKindNames[kind]
	if ok {
		return name
	}
	return "invalid"
}

func (kind NodeKind) String() string {
	return kind.Name()
}
