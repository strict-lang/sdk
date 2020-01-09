package analysis

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

const (
	MessageInvalidModuleImport    = "The imported modules name needs to be UpperCamelCase"
	MessageInvalidUnitName        = "The units name has to be lowerCamelCase"
	MessageInvalidDeclarationName = "Declared identifiers must be named lowerCamelCase"
	MessageImplicitParameterName  = "Parameters need explicit names if their type occurs more" +
		"than once in the parameter list"
)

// NamingCheck is traversing the tree and ensures that the names of all declared declaration
// follows the naming rules. As opposed to many languages, Strict opposes strong rules on the
// names of identifiers. Names may also influence the semantics.
type NamingCheck struct {
	recorder *diagnostic.Bag
	unit     *tree.TranslationUnit
	visitor  tree.Visitor
}

func NewNamingCheck(recorder *diagnostic.Bag, unit *tree.TranslationUnit) *NamingCheck {
	check := &NamingCheck{
		recorder: recorder,
		unit:     unit,
	}
	check.visitor = createVisitorForNamingCheck(check)
	return check
}

func createVisitorForNamingCheck(check *NamingCheck) tree.Visitor {
	visitor := tree.NewEmptyVisitor()
	visitor.ParameterVisitor = check.checkParameterNaming
	visitor.TranslationUnitVisitor = check.checkTranslationUnitNaming
	visitor.AssignStatementVisitor = check.checkFieldNaming
	visitor.MethodDeclarationVisitor = check.checkMethodNamingAndImplicitParameters
	visitor.ImportStatementVisitor = check.checkImportedModuleNaming
	visitor.ForEachLoopStatementVisitor = check.checkForEachLoopFieldNaming
	visitor.RangedLoopStatementVisitor = check.checkRangedLoopFieldNaming
	return visitor
}

func (check *NamingCheck) Run() {
	check.unit.AcceptRecursive(check.visitor)
}

// reportInvalidNode reports that the node has an invalid name.
func (check *NamingCheck) reportInvalidNode(node tree.Node, message string) {
	check.recorder.Record(diagnostic.RecordedEntry{
		Position: node.Locate(),
		UnitName: check.unit.Name,
		Kind:     &diagnostic.Error,
		Stage:    &diagnostic.SemanticAnalysis,
		Message:  message,
	})
}

// checkImportedModuleNaming ensures that the name of the imported module is upper camel case.
// Either by importing a file which starts with an upper case character or by having an
// alias that is upper camel case. Everything else results in a semantic error.
func (check *NamingCheck) checkImportedModuleNaming(statement *tree.ImportStatement) {
	if isUpperCamelCase(statement.ModuleName()) {
		check.reportInvalidNode(statement, MessageInvalidModuleImport)
	}
}

// checkRangedLoopFieldNaming ensures that the loops value fields naming is
// lowerCamelCase.
func (check *NamingCheck) checkRangedLoopFieldNaming(loop *tree.RangedLoopStatement) {
	if !isLowerCamelCase(loop.Field.Value) {
		check.reportInvalidNode(loop.Field, MessageInvalidDeclarationName)
	}
}

// checkForEachLoopFieldNaming ensures that the loops value fields  naming is
// lowerCamelCase.
func (check *NamingCheck) checkForEachLoopFieldNaming(loop *tree.ForEachLoopStatement) {
	if !isLowerCamelCase(loop.Field.Value) {
		check.reportInvalidNode(loop.Field, MessageInvalidDeclarationName)
	}
}

// checkTranslationUnitNaming ensures that the translation units name is a valid name for
// a Strict type, it has to be lowerCamelCase.
func (check *NamingCheck) checkTranslationUnitNaming(unit *tree.TranslationUnit) {
	if !isLowerCamelCase(unit.ToTypeName().NonGenericName()) {
		check.reportInvalidNode(unit, MessageInvalidUnitName)
	}
}

// checkFieldNaming ensures that all fields that are ever defined have an
// lowerCamelCase name.
func (check *NamingCheck) checkFieldNaming(method *tree.AssignStatement) {
	assignedField := method.Target
	identifier, isIdentifier := assignedField.(*tree.Identifier)
	if !isIdentifier {
		return
	}
	if !isLowerCamelCase(identifier.Value) {
		check.reportInvalidNode(identifier, MessageInvalidDeclarationName)
	}
}

// checkMethodNamingAndImplicitParameters ensures that a methods name is lowerCamelCase and that
// its parameters have explicit names if their type occurs more than once in the ParameterList.
// Meaning that when the ParameterList contains two numbers: '(number, number x)', both of the
// parameters need an explicit name: '(number x, number y)'.
func (check *NamingCheck) checkMethodNamingAndImplicitParameters(method *tree.MethodDeclaration) {
	if !isLowerCamelCase(method.Name.Value) {
		check.reportInvalidNode(method, MessageInvalidDeclarationName)
	}
	check.ensureExplicitParameterNamingOnDuplicateTypes(method.Parameters)
}

func (check *NamingCheck) ensureExplicitParameterNamingOnDuplicateTypes(parameters tree.ParameterList) {
	parameterTypeNames := map[string]bool{}
	for _, parameter := range parameters {
		if parameterTypeNames[parameter.Type.NonGenericName()] {
			check.reportInvalidNode(parameter, MessageImplicitParameterName)
		}
		parameterTypeNames[parameter.Type.NonGenericName()] = true
	}
}

// checkParameterNaming ensures that a parameter is named lowerCamelCase.
func (check *NamingCheck) checkParameterNaming(parameter *tree.Parameter) {
	if isLowerCamelCase(parameter.Name.Value) {
		return
	}
	check.reportInvalidNode(parameter, MessageInvalidDeclarationName)
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
