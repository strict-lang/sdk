package tree

// Visitor visits every node in the tree. The visitor-pattern is used to traverse
// the grammar-tree. Every tree-node has an 'Accept' method that lets the visitor
// visit itself and all of its children. In contrast to many visitor-pattern
// implementations, the visitor is not an interface. It has a lot of  methods
// and because of golang's type-system, every visitor-implementation would have
// to implement every method. Instead one can create a default-visitor, which
// only has noop-methods, and then set the visit-methods that should be custom.
type Visitor interface {
	VisitParameter(*Parameter)
	VisitIdentifier(*Identifier)
	VisitCallArgument(*CallArgument)
	VisitListTypeName(*ListTypeName)
	VisitTestStatement(*TestStatement)
	VisitStringLiteral(*StringLiteral)
	VisitNumberLiteral(*NumberLiteral)
	VisitCallExpression(*CallExpression)
	VisitEmptyStatement(*EmptyStatement)
	VisitYieldStatement(*YieldStatement)
	VisitBlockStatement(*BlockStatement)
	VisitAssertStatement(*AssertStatement)
	VisitUnaryExpression(*UnaryExpression)
	VisitImportStatement(*ImportStatement)
	VisitAssignStatement(*AssignStatement)
	VisitReturnStatement(*ReturnStatement)
	VisitTranslationUnit(*TranslationUnit)
	VisitCreateExpression(*CreateExpression)
	VisitInvalidStatement(*InvalidStatement)
	VisitFieldDeclaration(*FieldDeclaration)
	VisitGenericTypeName(*GenericTypeName)
	VisitConcreteTypeName(*ConcreteTypeName)
	VisitClassDeclaration(*ClassDeclaration)
	VisitBinaryExpression(*BinaryExpression)
	VisitMethodDeclaration(*MethodDeclaration)
	VisitRangedLoopStatement(*RangedLoopStatement)
	VisitExpressionStatement(*ExpressionStatement)
	VisitForEachLoopStatement(*ForEachLoopStatement)
	VisitConditionalStatement(*ConditionalStatement)
	VisitListSelectExpression(*ListSelectExpression)
	VisitConstructorDeclaration(*ConstructorDeclaration)
}

type DelegatingVisitor struct {
	ParameterVisitor              func(*Parameter)
	IdentifierVisitor             func(*Identifier)
	CallArgumentVisitor           func(*CallArgument)
	ListTypeNameVisitor           func(*ListTypeName)
	TestStatementVisitor          func(*TestStatement)
	StringLiteralVisitor          func(*StringLiteral)
	NumberLiteralVisitor          func(*NumberLiteral)
	CallExpressionVisitor         func(*CallExpression)
	EmptyStatementVisitor         func(*EmptyStatement)
	YieldStatementVisitor         func(*YieldStatement)
	BlockStatementVisitor         func(*BlockStatement)
	AssertStatementVisitor        func(*AssertStatement)
	UnaryExpressionVisitor        func(*UnaryExpression)
	ImportStatementVisitor        func(*ImportStatement)
	AssignStatementVisitor        func(*AssignStatement)
	ReturnStatementVisitor        func(*ReturnStatement)
	TranslationUnitVisitor        func(*TranslationUnit)
	CreateExpressionVisitor       func(*CreateExpression)
	InvalidStatementVisitor       func(*InvalidStatement)
	FieldDeclarationVisitor       func(*FieldDeclaration)
	GenericTypeNameVisitor        func(*GenericTypeName)
	ConcreteTypeNameVisitor       func(*ConcreteTypeName)
	ClassDeclarationVisitor       func(*ClassDeclaration)
	BinaryExpressionVisitor       func(*BinaryExpression)
	MethodDeclarationVisitor      func(*MethodDeclaration)
	RangedLoopStatementVisitor    func(*RangedLoopStatement)
	ExpressionStatementVisitor    func(*ExpressionStatement)
	ForEachLoopStatementVisitor   func(*ForEachLoopStatement)
	ConditionalStatementVisitor   func(*ConditionalStatement)
	ListSelectExpressionVisitor   func(*ListSelectExpression)
	ConstructorDeclarationVisitor func(*ConstructorDeclaration)
}

func NewEmptyVisitor() Visitor {
	return &DelegatingVisitor{
		ParameterVisitor:              func(*Parameter) {},
		IdentifierVisitor:             func(*Identifier) {},
		CallArgumentVisitor:           func(*CallArgument) {},
		ListTypeNameVisitor:           func(*ListTypeName) {},
		TestStatementVisitor:          func(*TestStatement) {},
		StringLiteralVisitor:          func(*StringLiteral) {},
		NumberLiteralVisitor:          func(*NumberLiteral) {},
		CallExpressionVisitor:         func(*CallExpression) {},
		EmptyStatementVisitor:         func(*EmptyStatement) {},
		YieldStatementVisitor:         func(*YieldStatement) {},
		BlockStatementVisitor:         func(*BlockStatement) {},
		AssertStatementVisitor:        func(*AssertStatement) {},
		UnaryExpressionVisitor:        func(*UnaryExpression) {},
		ImportStatementVisitor:        func(*ImportStatement) {},
		AssignStatementVisitor:        func(*AssignStatement) {},
		ReturnStatementVisitor:        func(*ReturnStatement) {},
		TranslationUnitVisitor:        func(*TranslationUnit) {},
		CreateExpressionVisitor:       func(*CreateExpression) {},
		InvalidStatementVisitor:       func(*InvalidStatement) {},
		FieldDeclarationVisitor:       func(*FieldDeclaration) {},
		GenericTypeNameVisitor:        func(*GenericTypeName) {},
		ConcreteTypeNameVisitor:       func(*ConcreteTypeName) {},
		ClassDeclarationVisitor:       func(*ClassDeclaration) {},
		BinaryExpressionVisitor:       func(*BinaryExpression) {},
		MethodDeclarationVisitor:      func(*MethodDeclaration) {},
		RangedLoopStatementVisitor:    func(*RangedLoopStatement) {},
		ExpressionStatementVisitor:    func(*ExpressionStatement) {},
		ForEachLoopStatementVisitor:   func(*ForEachLoopStatement) {},
		ConditionalStatementVisitor:   func(*ConditionalStatement) {},
		ListSelectExpressionVisitor:   func(*ListSelectExpression) {},
		ConstructorDeclarationVisitor: func(*ConstructorDeclaration) {},
	}
}
func (visitor *DelegatingVisitor) VisitParameter(node *Parameter) {
	visitor.ParameterVisitor(node)
}

func (visitor *DelegatingVisitor) VisitIdentifier(node *Identifier) {
	visitor.IdentifierVisitor(node)
}

func (visitor *DelegatingVisitor) VisitCallArgument(node *CallArgument) {
	visitor.CallArgumentVisitor(node)
}

func (visitor *DelegatingVisitor) VisitListTypeName(node *ListTypeName) {
	visitor.ListTypeNameVisitor(node)
}

func (visitor *DelegatingVisitor) VisitTestStatement(node *TestStatement) {
	visitor.TestStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitStringLiteral(node *StringLiteral) {
	visitor.StringLiteralVisitor(node)
}

func (visitor *DelegatingVisitor) VisitNumberLiteral(node *NumberLiteral) {
	visitor.NumberLiteralVisitor(node)
}

func (visitor *DelegatingVisitor) VisitCallExpression(node *CallExpression) {
	visitor.CallExpressionVisitor(node)
}

func (visitor *DelegatingVisitor) VisitEmptyStatement(node *EmptyStatement) {
	visitor.EmptyStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitYieldStatement(node *YieldStatement) {
	visitor.YieldStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitBlockStatement(node *BlockStatement) {
	visitor.BlockStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitAssertStatement(node *AssertStatement) {
	visitor.AssertStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitUnaryExpression(node *UnaryExpression) {
	visitor.UnaryExpressionVisitor(node)
}

func (visitor *DelegatingVisitor) VisitImportStatement(node *ImportStatement) {
	visitor.ImportStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitAssignStatement(node *AssignStatement) {
	visitor.AssignStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitReturnStatement(node *ReturnStatement) {
	visitor.ReturnStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitTranslationUnit(node *TranslationUnit) {
	visitor.TranslationUnitVisitor(node)
}

func (visitor *DelegatingVisitor) VisitCreateExpression(node *CreateExpression) {
	visitor.CreateExpressionVisitor(node)
}

func (visitor *DelegatingVisitor) VisitInvalidStatement(node *InvalidStatement) {
	visitor.InvalidStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitFieldDeclaration(node *FieldDeclaration) {
	visitor.FieldDeclarationVisitor(node)
}

func (visitor *DelegatingVisitor) VisitGenericTypeName(node *GenericTypeName) {
	visitor.GenericTypeNameVisitor(node)
}

func (visitor *DelegatingVisitor) VisitConcreteTypeName(node *ConcreteTypeName) {
	visitor.ConcreteTypeNameVisitor(node)
}

func (visitor *DelegatingVisitor) VisitClassDeclaration(node *ClassDeclaration) {
	visitor.ClassDeclarationVisitor(node)
}

func (visitor *DelegatingVisitor) VisitBinaryExpression(node *BinaryExpression) {
	visitor.BinaryExpressionVisitor(node)
}

func (visitor *DelegatingVisitor) VisitMethodDeclaration(node *MethodDeclaration) {
	visitor.MethodDeclarationVisitor(node)
}

func (visitor *DelegatingVisitor) VisitRangedLoopStatement(node *RangedLoopStatement) {
	visitor.RangedLoopStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitExpressionStatement(node *ExpressionStatement) {
	visitor.ExpressionStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitForEachLoopStatement(node *ForEachLoopStatement) {
	visitor.ForEachLoopStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitConditionalStatement(node *ConditionalStatement) {
	visitor.ConditionalStatementVisitor(node)
}

func (visitor *DelegatingVisitor) VisitListSelectExpression(node *ListSelectExpression) {
	visitor.ListSelectExpressionVisitor(node)
}

func (visitor *DelegatingVisitor) VisitConstructorDeclaration(node *ConstructorDeclaration) {
	visitor.ConstructorDeclarationVisitor(node)
}
