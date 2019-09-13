package syntaxtree

// Visitor visits every node in the syntaxtree. The visitor-pattern is used to traverse
// the syntax-tree. Every syntaxtree-node has an 'Accept' method that lets the visitor
// visit itself and all of its children. In contrast to many visitor-pattern
// implementations, the visitor is not an interface. It has a lot of  methods
// and because of golang's type-system, every visitor-implementation would have
// to implement every method. Instead one can create a default-visitor, which
// only has noop-methods, and then set the visit-methods that should be custom.
type Visitor struct {
	VisitParameter            func(*Parameter)
	VisitIdentifier           func(*Identifier)
	VisitCallArgument         func(*CallArgument)
	VisitListTypeName         func(*ListTypeName)
	VisitTestStatement        func(*TestStatement)
	VisitStringLiteral        func(*StringLiteral)
	VisitNumberLiteral        func(*NumberLiteral)
	VisitCallExpression       func(*CallExpression)
	VisitEmptyStatement       func(*EmptyStatement)
	VisitYieldStatement       func(*YieldStatement)
	VisitBlockStatement       func(*BlockStatement)
	VisitAssertStatement      func(*AssertStatement)
	VisitUnaryExpression      func(*UnaryExpression)
	VisitImportStatement      func(*ImportStatement)
	VisitAssignStatement      func(*AssignStatement)
	VisitReturnStatement      func(*ReturnStatement)
	VisitTranslationUnit      func(*TranslationUnit)
	VisitCreateExpression     func(*CreateExpression)
	VisitInvalidStatement     func(*InvalidStatement)
	VisitFieldDeclaration     func(*FieldDeclaration)
	VisitGenericTypeName      func(*GenericTypeName)
	VisitConcreteTypeName     func(*ConcreteTypeName)
	VisitClassDeclaration     func(*ClassDeclaration)
	VisitBinaryExpression     func(*BinaryExpression)
	VisitMethodDeclaration    func(*MethodDeclaration)
	VisitSelectorExpression   func(*SelectExpression)
	VisitIncrementStatement   func(*IncrementStatement)
	VisitDecrementStatement   func(*DecrementStatement)
	VisitRangedLoopStatement  func(*RangedLoopStatement)
	VisitExpressionStatement  func(*ExpressionStatement)
	VisitForEachLoopStatement func(*ForEachLoopStatement)
	VisitConditionalStatement func(*ConditionalStatement)
	VisitListSelectExpression func(*ListSelectExpression)
	VisitConstructorDeclaration func(*ConstructorDeclaration)
}

func NewEmptyVisitor() *Visitor {
	return &Visitor{
		VisitCallArgument:    func(*CallArgument) {},
		VisitParameter:       func(*Parameter) {},
		VisitCallExpression:  func(*CallExpression) {},
		VisitIdentifier:      func(*Identifier) {},
		VisitListTypeName:    func(*ListTypeName) {},
		VisitTestStatement:   func(*TestStatement) {},
		VisitStringLiteral:   func(*StringLiteral) {},
		VisitNumberLiteral:   func(*NumberLiteral) {},
		VisitYieldStatement:  func(*YieldStatement) {},
		VisitBlockStatement:  func(*BlockStatement) {},
		VisitAssertStatement: func(*AssertStatement) {},
		VisitUnaryExpression: func(*UnaryExpression) {},
		VisitCreateExpression:     func(*CreateExpression) {},
		VisitEmptyStatement:       func(*EmptyStatement) {},
		VisitImportStatement:      func(*ImportStatement) {},
		VisitReturnStatement:      func(*ReturnStatement) {},
		VisitTranslationUnit:      func(*TranslationUnit) {},
		VisitBinaryExpression:     func(*BinaryExpression) {},
		VisitAssignStatement:      func(*AssignStatement) {},
		VisitInvalidStatement:     func(*InvalidStatement) {},
		VisitFieldDeclaration:     func(*FieldDeclaration) {},
		VisitClassDeclaration:     func(*ClassDeclaration) {},
		VisitSelectorExpression:   func(*SelectExpression) {},
		VisitIncrementStatement:   func(*IncrementStatement) {},
		VisitDecrementStatement:   func(*DecrementStatement) {},
		VisitMethodDeclaration:    func(*MethodDeclaration) {},
		VisitGenericTypeName:      func(*GenericTypeName) {},
		VisitConcreteTypeName:     func(*ConcreteTypeName) {},
		VisitRangedLoopStatement:  func(*RangedLoopStatement) {},
		VisitExpressionStatement:  func(*ExpressionStatement) {},
		VisitForEachLoopStatement: func(*ForEachLoopStatement) {},
		VisitConditionalStatement: func(*ConditionalStatement) {},
		VisitListSelectExpression: func(*ListSelectExpression) {},
		VisitConstructorDeclaration: func(*ConstructorDeclaration) {},
	}
}
