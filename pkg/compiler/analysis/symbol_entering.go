package analysis

import (
	"errors"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/scope"
	"gitlab.com/strict-lang/sdk/pkg/compiler/typing"
)

// SymbolEnterPass enters symbols into the scope that they are defined in.
// It also ensures that there are no duplicates.
type SymbolEnterPass struct {
	diagnostics *diagnostic.Bag
	currentClass *tree.ClassDeclaration
	currentUnit *tree.TranslationUnit
}

func (pass *SymbolEnterPass) visitTranslationUnit(unit *tree.TranslationUnit) {
	pass.importIntoScope(unit, ensureScopeIsMutable(unit.Scope()))
}

func (pass *SymbolEnterPass) importIntoScope(
	unit *tree.TranslationUnit,
	targetScope scope.MutableScope) {

	for _, statement := range unit.Imports {
		pass.importNamespaceIntoScope(statement, targetScope)
	}
}

func (pass *SymbolEnterPass) importNamespaceIntoScope(
	statement *tree.ImportStatement,
	targetScope scope.MutableScope) {

	name := statement.Target.ToModuleName()
	imported, err := pass.captureExportedSymbolsOfNamespace(statement.Target)
	if err != nil {
		pass.reportImportError(statement.Target, err)
	}
	targetScope.Insert(&scope.Namespace{
		PackageName: name,
		Scope: imported,
	})
}

func (pass *SymbolEnterPass) reportImportError(
	namespace tree.ImportTarget, err error) {

}

func (pass *SymbolEnterPass) captureExportedSymbolsOfNamespace(
	namespace tree.ImportTarget) (scope.Scope, error) {

	return nil, errors.New("not implemented")
}

func (pass *SymbolEnterPass) visitClassDeclaration(
	declaration *tree.ClassDeclaration) {

	pass.currentClass = declaration
	pass.enterClassDeclaration(declaration)
}

func (pass *SymbolEnterPass) enterClassDeclaration(
	declaration *tree.ClassDeclaration) {

	name := declaration.Name
	surroundingScope := requireNearestMutableScope(declaration)
	if pass.ensureNameDoesNotExist(name, declaration, surroundingScope) {
		surroundingScope.Insert(pass.newClassSymbol(declaration))
	}
}

func (pass *SymbolEnterPass) newClassSymbol(
	declaration *tree.ClassDeclaration) *scope.Class {

	return &scope.Class{
		DeclarationName: declaration.Name,
		ActualClass:     declaration.NewActualClass(),
	}
}

func (pass *SymbolEnterPass) visitMethodDeclaration(
	declaration *tree.MethodDeclaration) {

	parameterSymbols := pass.enterMethodParameters(declaration)
	if symbol, ok := pass.enterMethodToSurroundingScope(declaration); ok {
		symbol.Parameters = parameterSymbols
	}
}

func (pass *SymbolEnterPass) enterMethodParameters(
	method *tree.MethodDeclaration) []*scope.Field {

	symbols := make([]*scope.Field, len(method.Parameters))
	methodScope := ensureScopeIsMutable(method.Scope())
	for index, parameter := range method.Parameters {
		symbols[index] = pass.enterMethodParameter(parameter, methodScope)
	}
	return symbols
}

func (pass *SymbolEnterPass) enterMethodParameter(
	parameter *tree.Parameter,
	methodScope scope.MutableScope) *scope.Field {

	symbol := pass.newFieldSymbolFromParameter(parameter, methodScope)
	if pass.ensureNameDoesNotExist(parameter.Name.Value, parameter, methodScope) {
		methodScope.Insert(symbol)
		return symbol
	}
	return nil
}

func (pass *SymbolEnterPass) newFieldSymbolFromParameter(
	parameter *tree.Parameter,
	surroundingScope scope.MutableScope) *scope.Field {

	class := pass.requireClass(parameter.Type, surroundingScope)
	return &scope.Field{
		Class:           class,
		DeclarationName: parameter.Name.Value,
		Kind:            scope.ParameterField,
	}
}

func (pass *SymbolEnterPass) enterMethodToSurroundingScope(
	method *tree.MethodDeclaration) (*scope.Method, bool) {

	surroundingScope := requireNearestMutableScope(method)
	if pass.ensureNameDoesNotExist(method.Name.Value, method, surroundingScope) {
		symbol := pass.newMethodSymbol(method, surroundingScope)
		surroundingScope.Insert(symbol)
		return symbol, true
	}
	return nil, false
}

func (pass *SymbolEnterPass) newMethodSymbol(
	method *tree.MethodDeclaration,
	surroundingScope scope.MutableScope) *scope.Method {

	class := pass.requireClass(method.Type, surroundingScope)
	return &scope.Method{
		ReturnType: class,
	}
}

// requireClass is used to resolve classes of certain declarations. While this
// pass is mainly inserting stuff into the scope, it has to get a reference of
// the classes for return types and parameter/field types. If no class with the
// given name exists, an error is reported and a "class replacement" is created
// in the targetScope. This replacement is an empty class with the required name.
// By doing this, the name resolution can continue and more semantic errors may
// be provided to the user, even if the method declaration itself is invalid.
func (pass *SymbolEnterPass) requireClass(
	name tree.TypeName, targetScope scope.MutableScope) *scope.Class {

	returnTypePoint := scope.NewReferencePoint(name.NonGenericName())
	if class, ok := scope.LookupClass(targetScope, returnTypePoint); ok {
		return class
	}
	pass.reportMissingClass(name)
	return pass.createClassReplacementInScope(name, targetScope)
}

func (pass *SymbolEnterPass) reportMissingClass(name tree.TypeName) {

}

func (pass *SymbolEnterPass) createClassReplacementInScope(
	name tree.TypeName,
	targetScope scope.MutableScope) *scope.Class {

	symbol := pass.createClassReplacement(name)
	targetScope.Insert(symbol)
	return symbol
}

func (pass *SymbolEnterPass) createClassReplacement(
	name tree.TypeName) *scope.Class {

	return &scope.Class{
		DeclarationName: name.NonGenericName(),
		ActualClass:     typing.NewEmptyClass(name.NonGenericName()),
	}
}

func (pass *SymbolEnterPass) ensureNameDoesNotExist(
	name string, node tree.Node, surroundingScope scope.Scope) bool {

	point := scope.NewReferencePoint(name)
	if entries := surroundingScope.Lookup(point); !entries.IsEmpty() {
		existingSymbol := entries.First().Symbol
		pass.reportNameCollision(name, node, existingSymbol)
		return false
	}
	return true
}

func (pass *SymbolEnterPass) reportNameCollision(
	name string,
	node tree.Node,
	existingSymbol scope.Symbol) {

}
