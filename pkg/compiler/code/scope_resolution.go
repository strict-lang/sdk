package code

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
	diagnostics *diagnostic.Bag
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
	unit *tree.TranslationUnit) { }

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
	surroundingScope, exists := tree.ResolveNearestScope(node)
	if !exists {
		log.Fatalf("surrounding scope does not exists: %v", node)
	}
	return surroundingScope
}
