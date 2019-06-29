package ast

// Visitor visits every node in the ast. The visitor-pattern is used to traverse
// the syntax-tree. Every ast-node has an 'Accept' method that lets the visitor
// visit itself and all of its children.
type Visitor struct {
	// VisitMethodCall visits a MethodCall node.
	VisitMethodCall func(*MethodCall)

	// VisitIdentifier visits an Identifier node.
	VisitIdentifier func(*Identifier)

	// VisitStringLiteral visits a StringLiteral node.
	VisitStringLiteral func(*StringLiteral)

	// VisitNumberLiteral visits a NumberLiteral node.
	VisitNumberLiteral func(*NumberLiteral)

	// VisitBlockStatement visits a BlockStatement node.
	VisitBlockStatement func(*BlockStatement)

	// VisitUnaryExpression visits an UnaryExpression node.
	VisitUnaryExpression func(*UnaryExpression)

	// VisitTranslationUnit vitits a TranslationUnit node.
	VisitTranslationUnit func(*TranslationUnit)

	// VisitForLoopStatement visits a ForLoopStatement node.
	VisitForLoopStatement func(*ForLoopStatement)

	// VisitBinaryExpression visits an BinaryExpression node.
	VisitBinaryExpression func(*BinaryExpression)

	// VisitExpressionStatement visits an ExpressionStatement node.
	VisitExpressionStatement func(*ExpressionStatement)

	// VisitConditionalStatement visits a ConditionalStatement node.
	VisitConditionalStatement func(*ConditionalStatement)
}
