package ast

// Visitor visits every node in the ast. The visitor-pattern is used to traverse
// the syntax-tree. Every ast-node has an 'Accept' method that lets the visitor
//visit itself and all of its children. In contrast to many visitor-pattern
// implementations, the visitor is not an interface. It has a lot of  methods
// and because of golang's type-system, every visitor-implementation would have
// to implement every method. Instead one can create a default-visitor, which
// only has noop-methods, and then set the visit-methods that should be custom.
type Visitor struct {
	// VisitType visits a Type node.
	VisitType func(*Type)

	// VisitMember visits a Member node.
	VisitMember func(*Member)

	// VisitMethod visits a Method node.
	VisitMethod func(*Method)

	// VisitParameter visits a Parameter node.
	VisitParameter func(*Parameter)

	// VisitMethodCall visits a MethodCall node.
	VisitMethodCall func(*MethodCall)

	// VisitIdentifier visits an Identifier node.
	VisitIdentifier func(*Identifier)

	// VisitStringLiteral visits a StringLiteral node.
	VisitStringLiteral func(*StringLiteral)

	// VisitNumberLiteral visits a NumberLiteral node.
	VisitNumberLiteral func(*NumberLiteral)

	// VisitEmptyStatement visits a EmptyStatement node.
	VisitEmptyStatement func(*EmptyStatement)

	// VisitYieldStatement visits a YieldStatement node.
	VisitYieldStatement func(*YieldStatement)

	// VisitBlockStatement visits a BlockStatement node.
	VisitBlockStatement func(*BlockStatement)

	// VisitUnaryExpression visits an UnaryExpression node.
	VisitUnaryExpression func(*UnaryExpression)

	// VisitAssignStatement visits an AssignStatement node.
	VisitAssignStatement func(*AssignStatement)

	// VisitReturnStatement visits a ReturnStatement node.
	VisitReturnStatement func(*ReturnStatement)

	// VisitTranslationUnit visits a TranslationUnit node.
	VisitTranslationUnit func(*TranslationUnit)

	// VisitInvalidStatement visits an InvalidStatement node.
	VisitInvalidStatement func (*InvalidStatement)

	// VisitBinaryExpression visits an BinaryExpression node.
	VisitBinaryExpression func(*BinaryExpression)

	// VisitIncrementStatement visits a IncrementStatement node.
	VisitIncrementStatement func(*IncrementStatement)

	// VisitDecrementStatement visits a DecrementStatement node.
	VisitDecrementStatement func(*DecrementStatement)

	// VisitFromToLoopStatement visits a FromToLoopStatement node.
	VisitFromToLoopStatement func(*FromToLoopStatement)

	// VisitExpressionStatement visits an ExpressionStatement node.
	VisitExpressionStatement func(*ExpressionStatement)

	// VisitForeachLoopStatement visits a ForeachLoopStatement node.
	VisitForeachLoopStatement func(*ForeachLoopStatement)

	// VisitConditionalStatement visits a ConditionalStatement node.
	VisitConditionalStatement func(*ConditionalStatement)

}

func NewEmptyVisitor() *Visitor {
	return &Visitor{
		VisitType:                   func(*Type) {},
		VisitMember:                 func(*Member) {},
		VisitMethod:                 func(*Method) {},
		VisitParameter:              func(*Parameter) {},
		VisitMethodCall:             func(*MethodCall) {},
		VisitIdentifier:             func(*Identifier) {},
		VisitStringLiteral:          func(*StringLiteral) {},
		VisitNumberLiteral:          func(*NumberLiteral) {},
		VisitYieldStatement:         func(*YieldStatement) {},
		VisitBlockStatement:         func(*BlockStatement) {},
		VisitUnaryExpression:        func(*UnaryExpression) {},
		VisitEmptyStatement: 				 func(*EmptyStatement) {},
		VisitReturnStatement:        func(*ReturnStatement) {},
		VisitTranslationUnit:        func(*TranslationUnit) {},
		VisitBinaryExpression:       func(*BinaryExpression) {},
		VisitAssignStatement:        func(*AssignStatement) {},
		VisitInvalidStatement:       func(*InvalidStatement) {},
		VisitIncrementStatement:  	 func(*IncrementStatement) {},
		VisitDecrementStatement:     func(*DecrementStatement) {},
		VisitFromToLoopStatement:    func(*FromToLoopStatement) {},
		VisitExpressionStatement:    func(*ExpressionStatement) {},
		VisitForeachLoopStatement:   func(*ForeachLoopStatement) {},
		VisitConditionalStatement:   func(*ConditionalStatement) {},
	}
}
