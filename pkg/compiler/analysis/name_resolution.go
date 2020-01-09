package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
)

type NameResolutionPass struct {
}

type chainResolution struct {
	currentScope scope.Scope
	currentNode tree.Node
}

func (resolution *chainResolution) bindAll() {

}

func (resolution *NameResolutionPass) visitFieldSelect(
	expression *tree.FieldSelectExpression) {

	chainResolution := &chainResolution{
		currentScope: requireNearestScope(expression),
		currentNode:  expression,
	}
	chainResolution.bindAll()
}

func (resolution *NameResolutionPass) visitIdentifier(
	identifier *tree.Identifier) {

	if !identifier.IsBound() && !identifier.IsPartOfDeclaration() {
		resolution.bindIdentifier(identifier)
	}
}

func (resolution *NameResolutionPass) bindIdentifier(
	identifier *tree.Identifier) {

	surroundingScope := requireNearestScope(identifier)
	point := identifier.ReferencePoint()
	if entries := surroundingScope.Lookup(point); !entries.IsEmpty() {
		symbol := entries.First().Symbol
		identifier.Bind(symbol)
	} else {

	}
}

func (resolution *NameResolutionPass) reportUnknownIdentifier(
	identifier *tree.Identifier) {


}