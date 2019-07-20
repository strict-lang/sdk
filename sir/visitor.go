package sir

type Visitor struct {
	VisitUnit              func(*Unit)
	VisitModule            func(*Module)
	VisitMethodParameter   func(*MethodParameter)
	VisitMethodDeclaration func(*MethodDeclaration)
}
