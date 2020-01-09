package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
	"log"
)

func (resolution *ScopeResolution) createScopeResolutionVisitor() tree.Visitor {
	return &tree.DelegatingVisitor{
		BlockStatementVisitor:         resolution.createBlockStatementScope,
		TranslationUnitVisitor:        resolution.createTranslationUnitScope,
		ClassDeclarationVisitor:       resolution.createClassDeclarationScope,
		MethodDeclarationVisitor:      resolution.createMethodDeclarationScope,
		ConstructorDeclarationVisitor: resolution.createConstructorDeclarationScope,
	}
}

type ScopeResolution struct {
	diagnostics  *diagnostic.Bag
	localIdCount int
}

func (*ScopeResolution) createMethodDeclarationScope(method *tree.MethodDeclaration) {
	surroundingScope := requireNearestScope(method)
	localScope := scope.NewLocalScope(
		scope.Id(method.Name.Value),
		method.Region,
		surroundingScope)
	method.UpdateScope(localScope)
}

func (resolution *ScopeResolution) createBlockStatementScope(
	block *tree.StatementBlock) {

	surroundingScope := requireNearestScope(block)
	localScope := scope.NewLocalScope(
		resolution.nextLocalIdSuffix(),
		block.Region,
		surroundingScope)
	block.UpdateScope(localScope)
}

func (resolution *ScopeResolution) createTranslationUnitScope(
	unit *tree.TranslationUnit) {}

func (resolution *ScopeResolution) createClassDeclarationScope(
	class *tree.ClassDeclaration) {

	surroundingScope := requireNearestScope(class)
	localScope := scope.NewLocalScope(
		resolution.nextLocalIdSuffix(),
		class.Region,
		surroundingScope)
	class.UpdateScope(localScope)
}

func (resolution *ScopeResolution) createConstructorDeclarationScope(
	constructor *tree.ConstructorDeclaration) {

	surroundingScope := requireNearestScope(constructor)
	localScope := scope.NewLocalScope(
		resolution.nextLocalIdSuffix(),
		constructor.Region,
		surroundingScope)
	constructor.UpdateScope(localScope)
}

func (resolution *ScopeResolution) nextLocalIdSuffix() scope.Id {
	resolution.localIdCount++
	return scope.Id(resolution.localIdCount)
}

func requireNearestScope(node tree.Node) scope.Scope {
	if surroundingScope, ok := tree.ResolveNearestScope(node); ok {
		return surroundingScope
	}
	log.Fatalf("surrounding scope does not exist: %v", node)
	return nil
}

func requireNearestMutableScope(node tree.Node) scope.MutableScope {
	if surroundingScope, ok := tree.ResolveNearestMutableScope(node); ok {
		return surroundingScope
	}
	log.Fatalf("surrounding mutable scope does not exist: %v", node)
	return nil
}

func ensureScopeIsMutable(anyScope scope.Scope) scope.MutableScope {
	if mutable, ok := anyScope.(scope.MutableScope); ok {
		return mutable
	}
	log.Fatal("scope is not mutable")
	return nil
}
