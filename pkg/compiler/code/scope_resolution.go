package code

import "gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"

func runScopeResolution(unit *tree.TranslationUnit) {

}

func createScopeResolutionVisitor() tree.Visitor {
	return &tree.DelegatingVisitor{
		BlockStatementVisitor:         nil,
		TranslationUnitVisitor:        nil,
		ClassDeclarationVisitor:       nil,
		MethodDeclarationVisitor:      createMethodDeclarationScope,
		ConstructorDeclarationVisitor: nil,
	}
}

func createMethodDeclarationScope(method *tree.MethodDeclaration) {
}
