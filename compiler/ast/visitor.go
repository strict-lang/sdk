package ast

// Visitor visits every node in the ast. The visitor-pattern is used to traverse
// the syntax-tree. Every ast-node has an 'Accept' method that lets the visitor
//visit itself and all of its children. In contrast to many visitor-pattern
// implementations, the visitor is not an interface. It has a lot of  methods
// and because of golang's type-system, every visitor-implementation would have
// to implement every method. Instead one can create a default-visitor, which
// only has noop-methods, and then set the visit-methods that should be custom.
type Visitor struct {
	VisitType                 func(*Type)

	VisitParameter            func(*Parameter)
	VisitMethodCall           func(*MethodCall)
	VisitIdentifier           func(*Identifier)
	VisitStringLiteral        func(*StringLiteral)
	VisitNumberLiteral        func(*NumberLiteral)
	VisitEmptyStatement       func(*EmptyStatement)
	VisitYieldStatement       func(*YieldStatement)
	VisitBlockStatement       func(*BlockStatement)
	VisitUnaryExpression      func(*UnaryExpression)
	VisitImportStatement      func(*ImportStatement)
	VisitAssignStatement      func(*AssignStatement)
	VisitReturnStatement      func(*ReturnStatement)
	VisitTranslationUnit      func(*TranslationUnit)
	VisitInvalidStatement     func(*InvalidStatement)
	VisitBinaryExpression     func(*BinaryExpression)
	VisitMethodDeclaration    func(*MethodDeclaration)
	VisitSelectorExpression   func(*SelectorExpression)
	VisitIncrementStatement   func(*IncrementStatement)
	VisitDecrementStatement   func(*DecrementStatement)
	VisitRangedLoopStatement  func(*RangedLoopStatement)
	VisitExpressionStatement  func(*ExpressionStatement)
	VisitForEachLoopStatement func(*ForEachLoopStatement)
	VisitConditionalStatement func(*ConditionalStatement)
}

func NewEmptyVisitor() *Visitor {
	return &Visitor{
		VisitParameter:            func(*Parameter) {},
		VisitMethodCall:           func(*MethodCall) {},
		VisitIdentifier:           func(*Identifier) {},
		VisitStringLiteral:        func(*StringLiteral) {},
		VisitNumberLiteral:        func(*NumberLiteral) {},
		VisitYieldStatement:       func(*YieldStatement) {},
		VisitBlockStatement:       func(*BlockStatement) {},
		VisitUnaryExpression:      func(*UnaryExpression) {},
		VisitEmptyStatement:       func(*EmptyStatement) {},
		VisitImportStatement:      func(*ImportStatement) {},
		VisitReturnStatement:      func(*ReturnStatement) {},
		VisitTranslationUnit:      func(*TranslationUnit) {},
		VisitBinaryExpression:     func(*BinaryExpression) {},
		VisitAssignStatement:      func(*AssignStatement) {},
		VisitInvalidStatement:     func(*InvalidStatement) {},
		VisitSelectorExpression:   func(*SelectorExpression) {},
		VisitIncrementStatement:   func(*IncrementStatement) {},
		VisitDecrementStatement:   func(*DecrementStatement) {},
		VisitMethodDeclaration:               func(*MethodDeclaration) {},
		VisitRangedLoopStatement:  func(*RangedLoopStatement) {},
		VisitExpressionStatement:  func(*ExpressionStatement) {},
		VisitForEachLoopStatement: func(*ForEachLoopStatement) {},
		VisitConditionalStatement: func(*ConditionalStatement) {},
		VisitSharedVariableDeclaration: func(*SharedVariableDeclaration) {},
	}
}
