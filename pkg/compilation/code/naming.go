package code

import (
	diagnostic2 "gitlab.com/strict-lang/sdk/pkg/compilation/diagnostic"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

const (
	MessageInvalidModuleImport    = "The imported modules name needs to be UpperCamelCase"
	MessageInvalidUnitName        = "The units name has to be lowerCamelCase"
	MessageInvalidDeclarationName = "Declared identifiers must be named lowerCamelCase"
	MessageImplicitParameterName  = "Parameters need explicit names if their type occurs more" +
		"than once in the parameter list"
)

// NamingCheck is traversing the syntaxtree and ensures that the names of all declared declaration
// follows the naming rules. As opposed to many languages, Strict opposes strong rules on the
// names of identifiers. Names may also influence the semantics.
type NamingCheck struct {
	recorder *diagnostic2.Bag
	unit     *syntaxtree2.TranslationUnit
	visitor  *syntaxtree2.Visitor
}

func NewNamingChecker(recorder *diagnostic2.Bag, unit *syntaxtree2.TranslationUnit) *NamingCheck {
	check := &NamingCheck{
		recorder: recorder,
		unit:     unit,
	}
	check.visitor = syntaxtree2.NewEmptyVisitor()
	check.visitor.VisitParameter = check.CheckParameterNaming
	check.visitor.VisitTranslationUnit = check.CheckTranslationUnitNaming
	check.visitor.VisitAssignStatement = check.CheckFieldNaming
	check.visitor.VisitMethodDeclaration = check.CheckMethodNamingAndImplicitParameters
	check.visitor.VisitImportStatement = check.CheckImportedModuleNaming
	check.visitor.VisitForEachLoopStatement = check.CheckForEachLoopFieldNaming
	check.visitor.VisitRangedLoopStatement = check.CheckRangedLoopFieldNaming
	return check
}

func (check *NamingCheck) Run() {
	check.unit.AcceptRecursive(check.visitor)
}

// reportInvalidNode reports that the node has an invalid name.
func (check *NamingCheck) reportInvalidNode(node syntaxtree2.Node, message string) {
	check.recorder.Record(diagnostic2.RecordedEntry{
		Position: node.Position(),
		UnitName: check.unit.Name,
		Kind:     &diagnostic2.Error,
		Stage:    &diagnostic2.SemanticAnalysis,
		Message:  message,
	})
}

// CheckImportedModuleNaming ensures that the name of the imported module is upper camel case.
// Either by importing a file which starts with an upper case character or by having an
// alias that is upper camel case. Everything else results in a semantic error.
func (check *NamingCheck) CheckImportedModuleNaming(statement *syntaxtree2.ImportStatement) {
	if isUpperCamelCase(statement.ModuleName()) {
		check.reportInvalidNode(statement, MessageInvalidModuleImport)
	}
}

// CheckRangedLoopFieldNaming ensures that the loops value fields naming is lowerCamelCase.
func (check *NamingCheck) CheckRangedLoopFieldNaming(loop *syntaxtree2.RangedLoopStatement) {
	if !isLowerCamelCase(loop.ValueField.Value) {
		check.reportInvalidNode(loop.ValueField, MessageInvalidDeclarationName)
	}
}

// CheckForEachLoopFieldNaming ensures that the loops value fields  naming is lowerCamelCase.
func (check *NamingCheck) CheckForEachLoopFieldNaming(loop *syntaxtree2.ForEachLoopStatement) {
	if !isLowerCamelCase(loop.Field.Value) {
		check.reportInvalidNode(loop.Field, MessageInvalidDeclarationName)
	}
}

// CheckTranslationUnitNaming ensures that the translation units name is a valid name for
// a Strict type, it has to be lowerCamelCase.
func (check *NamingCheck) CheckTranslationUnitNaming(unit *syntaxtree2.TranslationUnit) {
	if !isLowerCamelCase(unit.ToTypeName().NonGenericName()) {
		check.reportInvalidNode(unit, MessageInvalidUnitName)
	}
}

// CheckFieldNaming ensures that all fields that are ever defined have an lowerCamelCase name.
func (check *NamingCheck) CheckFieldNaming(method *syntaxtree2.AssignStatement) {
	assignedField := method.Target
	identifier, isIdentifier := assignedField.(*syntaxtree2.Identifier)
	if !isIdentifier {
		return
	}
	if !isLowerCamelCase(identifier.Value) {
		check.reportInvalidNode(identifier, MessageInvalidDeclarationName)
	}
}

// CheckMethodNamingAndImplicitParameters ensures that a methods name is lowerCamelCase and that
// its parameters have explicit names if their type occurs more than once in the ParameterList.
// Meaning that when the ParameterList contains two numbers: '(number, number x)', both of the
// parameters need an explicit name: '(number x, number y)'.
func (check *NamingCheck) CheckMethodNamingAndImplicitParameters(method *syntaxtree2.MethodDeclaration) {
	if !isLowerCamelCase(method.Name.Value) {
		check.reportInvalidNode(method, MessageInvalidDeclarationName)
	}
	check.ensureExplicitParameterNamingOnDuplicateTypes(method.Parameters)
}

func (check *NamingCheck) ensureExplicitParameterNamingOnDuplicateTypes(parameters syntaxtree2.ParameterList) {
	parameterTypeNames := map[string]bool{}
	for _, parameter := range parameters {
		if parameterTypeNames[parameter.Type.NonGenericName()] {
			check.reportInvalidNode(parameter, MessageImplicitParameterName)
		}
		parameterTypeNames[parameter.Type.NonGenericName()] = true
	}
}

// CheckParameterNaming ensures that a parameter is named lowerCamelCase.
func (check *NamingCheck) CheckParameterNaming(parameter *syntaxtree2.Parameter) {
	if isLowerCamelCase(parameter.Name.Value) {
		return
	}
	if parameter.IsNamedAfterType() {
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