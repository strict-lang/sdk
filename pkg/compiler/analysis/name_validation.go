package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	passes "gitlab.com/strict-lang/sdk/pkg/compiler/pass"
)

const (
	MessageInvalidModuleImport    = "The imported modules name needs to be UpperCamelCase"
	MessageInvalidUnitName        = "The units name has to be lowerCamelCase"
	MessageInvalidDeclarationName = "Declared identifiers must be named lowerCamelCase"
	MessageImplicitParameterName  = "Parameters need explicit names if their type occurs more" +
		"than once in the parameter list"
)

const NamingCheckPassId = "NamingCheckPass"

// NamingCheckPass is traversing the tree and ensures that the names of all declared declaration
// follows the naming rules. As opposed to many languages, Strict opposes strong rules on the
// names of identifiers. Names may also influence the semantics.
type NamingCheckPass struct {
	recorder *diagnostic.Bag
	unit     *tree.TranslationUnit
	visitor  tree.Visitor
}

func (pass *NamingCheckPass) Run(context *passes.Context) {
	visitor := pass.createVisitor()
	context.Unit.AcceptRecursive(visitor)
}

func (pass *NamingCheckPass) Dependencies(*isolate.Isolate) passes.Set {
	return passes.EmptySet
}

func (pass *NamingCheckPass) createVisitor() tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.ParameterVisitor = pass.checkParameterNaming
	visitor.TranslationUnitVisitor = pass.checkTranslationUnitNaming
	visitor.AssignStatementVisitor = pass.checkFieldNaming
	visitor.MethodDeclarationVisitor = pass.checkMethodNamingAndImplicitParameters
	visitor.ImportStatementVisitor = pass.checkImportedModuleNaming
	visitor.ForEachLoopStatementVisitor = pass.checkForEachLoopFieldNaming
	visitor.RangedLoopStatementVisitor = pass.checkRangedLoopFieldNaming
	return visitor
}

// reportInvalidNode reports that the node has an invalid name.
func (pass *NamingCheckPass) reportInvalidNode(node tree.Node, message string) {
	pass.recorder.Record(diagnostic.RecordedEntry{
		Position: node.Locate(),
		UnitName: pass.unit.Name,
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SemanticAnalysis,
		Message:  message,
	})
}

// checkImportedModuleNaming ensures that the name of the imported module is upper camel case.
// Either by importing a file which starts with an upper case character or by having an
// alias that is upper camel case. Everything else results in a semantic error.
func (pass *NamingCheckPass) checkImportedModuleNaming(statement *tree.ImportStatement) {
	if isUpperCamelCase(statement.ModuleName()) {
		pass.reportInvalidNode(statement, MessageInvalidModuleImport)
	}
}

// checkRangedLoopFieldNaming ensures that the loops value fields naming is
// lowerCamelCase.
func (pass *NamingCheckPass) checkRangedLoopFieldNaming(loop *tree.RangedLoopStatement) {
	if !isLowerCamelCase(loop.Field.Value) {
		pass.reportInvalidNode(loop.Field, MessageInvalidDeclarationName)
	}
}

// checkForEachLoopFieldNaming ensures that the loops value fields  naming is
// lowerCamelCase.
func (pass *NamingCheckPass) checkForEachLoopFieldNaming(loop *tree.ForEachLoopStatement) {
	if !isLowerCamelCase(loop.Field.Value) {
		pass.reportInvalidNode(loop.Field, MessageInvalidDeclarationName)
	}
}

// checkTranslationUnitNaming ensures that the translation units name is a valid name for
// a Strict type, it has to be lowerCamelCase.
func (pass *NamingCheckPass) checkTranslationUnitNaming(unit *tree.TranslationUnit) {
	if !isLowerCamelCase(unit.ToTypeName().NonGenericName()) {
		pass.reportInvalidNode(unit, MessageInvalidUnitName)
	}
}

// checkFieldNaming ensures that all fields that are ever defined have an
// lowerCamelCase name.
func (pass *NamingCheckPass) checkFieldNaming(method *tree.AssignStatement) {
	assignedField := method.Target
	identifier, isIdentifier := assignedField.(*tree.Identifier)
	if !isIdentifier {
		return
	}
	if !isLowerCamelCase(identifier.Value) {
		pass.reportInvalidNode(identifier, MessageInvalidDeclarationName)
	}
}

// checkMethodNamingAndImplicitParameters ensures that a methods name is lowerCamelCase and that
// its parameters have explicit names if their type occurs more than once in the ParameterList.
// Meaning that when the ParameterList contains two numbers: '(number, number x)', both of the
// parameters need an explicit name: '(number x, number y)'.
func (pass *NamingCheckPass) checkMethodNamingAndImplicitParameters(method *tree.MethodDeclaration) {
	if !isLowerCamelCase(method.Name.Value) {
		pass.reportInvalidNode(method, MessageInvalidDeclarationName)
	}
	pass.ensureExplicitParameterNamingOnDuplicateTypes(method.Parameters)
}

func (pass *NamingCheckPass) ensureExplicitParameterNamingOnDuplicateTypes(parameters tree.ParameterList) {
	parameterTypeNames := map[string]bool{}
	for _, parameter := range parameters {
		if parameterTypeNames[parameter.Type.NonGenericName()] {
			pass.reportInvalidNode(parameter, MessageImplicitParameterName)
		}
		parameterTypeNames[parameter.Type.NonGenericName()] = true
	}
}

// checkParameterNaming ensures that a parameter is named lowerCamelCase.
func (pass *NamingCheckPass) checkParameterNaming(parameter *tree.Parameter) {
	if isLowerCamelCase(parameter.Name.Value) {
		return
	}
	pass.reportInvalidNode(parameter, MessageInvalidDeclarationName)
}

func isCharLowerCase(char uint8) bool {
	return char >= 'a' && char <= 'z'
}

func isCharUpperCase(char uint8) bool {
	return char >= 'A' && char <= 'Z'
}

func isLowerCamelCase(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if !isCharLowerCase(identifier[0]) {
		return false
	}
	if len(identifier) == 1 {
		return true
	}
	return isValidLowerCamelCaseTail(identifier[1:])
}

func isValidLowerCamelCaseTail(identifier string) bool {
	var wasLastUpperCase = false
	for index := 0; index < len(identifier); index++ {
		switch char := identifier[index]; {
		case isCharLowerCase(char):
			wasLastUpperCase = false
			continue
		case isCharUpperCase(char):
			if wasLastUpperCase {
				return false
			}
			wasLastUpperCase = true
			continue
		}
	}
	return true
}

func isUpperCamelCase(identifier string) bool {
	if len(identifier) == 0 {
		return false
	}
	if !isCharUpperCase(identifier[0]) {
		return false
	}
	if len(identifier) == 1 {
		return true
	}
	// Reusing the lower camel case tail function, since it
	// prevents two consecutive upper case characters while
	// still allowing upper case characters to occur.
	return isValidLowerCamelCaseTail(identifier[1:])
}
