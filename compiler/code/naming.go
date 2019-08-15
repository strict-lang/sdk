package code

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/diagnostic"
)

const (
	MessageInvalidModuleImport = "The imported modules name needs to be UpperCamelCase"
	MessageInvalidUnitName = "The units name has to be lowerCamelCase"
	MessageInvalidDeclarationName = "Declared identifiers must be named lowerCamelCase"
	MessageImplicitParameterName = "Parameters need explicit names if their type occurs more" +
		"than once in the parameter list"
)

// NamingChecker is traversing the ast and ensures that the names of all declared declaration
// follows the naming rules. As opposed to many languages, strict opposes strong rules on the
// names of identifiers. Names may also influence the semantics.
type NamingChecker struct {
	recorder *diagnostic.Recorder
	unit *ast.TranslationUnit
	visitor *ast.Visitor
}

func NewNamingChecker(recorder *diagnostic.Recorder, unit *ast.TranslationUnit) *NamingChecker {
	checker := &NamingChecker{
		recorder: recorder,
		unit: unit,
	}
	checker.visitor = ast.NewEmptyVisitor()
	checker.visitor.VisitParameter = checker.CheckParameterNaming
	checker.visitor.VisitTranslationUnit = checker.CheckTranslationUnitNaming
	checker.visitor.VisitAssignStatement = checker.CheckFieldNaming
	checker.visitor.VisitMethodDeclaration = checker.CheckMethodNamingAndImplicitParameters
	checker.visitor.VisitImportStatement = checker.CheckImportedModuleNaming
	checker.visitor.VisitForEachLoopStatement = checker.CheckForEachLoopFieldNaming
	checker.visitor.VisitRangedLoopStatement = checker.CheckRangedLoopFieldNaming
	return checker
}

func (checker *NamingChecker) Run() {
	checker.unit.AcceptAll(checker.visitor)
}

// reportInvalidNode reports that the node has an invalid name.
func (checker *NamingChecker) reportInvalidNode(node ast.Node, message string) {
	checker.recorder.Record(diagnostic.RecordedEntry{
		Position: node.Position(),
		UnitName: checker.unit.Name(),
		Kind: &diagnostic.Error,
		Stage: &diagnostic.SemanticAnalysis,
		Message: message,
	})
}

// CheckImportedModuleNaming ensures that the name of the imported module is upper camel case.
// Either by importing a file which starts with an upper case character or by having an
// alias that is upper camel case. Everything else results in a semantic error.
func (checker *NamingChecker) CheckImportedModuleNaming(statement *ast.ImportStatement) {
	if isUpperCamelCase(statement.ModuleName()) {
		checker.reportInvalidNode(statement, MessageInvalidModuleImport)
	}
}

// CheckRangedLoopFieldNaming ensures that the loops value fields naming is lowerCamelCase.
func (checker *NamingChecker) CheckRangedLoopFieldNaming(loop *ast.RangedLoopStatement) {
	if !isLowerCamelCase(loop.ValueField.Value) {
		checker.reportInvalidNode(loop.ValueField, MessageInvalidDeclarationName)
	}
}

// CheckForEachLoopFieldNaming ensures that the loops value fields  naming is lowerCamelCase.
func (checker *NamingChecker) CheckForEachLoopFieldNaming(loop *ast.ForEachLoopStatement) {
	if !isLowerCamelCase(loop.Field.Value) {
		checker.reportInvalidNode(loop.Field, MessageInvalidDeclarationName)
	}
}

// CheckTranslationUnitNaming ensures that the translation units name is a valid name for
// a strict type, it has to be lowerCamelCase.
func (checker *NamingChecker) CheckTranslationUnitNaming(unit *ast.TranslationUnit) {
	if !isLowerCamelCase(unit.ToTypeName().NonGenericName()) {
		checker.reportInvalidNode(unit, MessageInvalidUnitName)
	}
}

// CheckFieldNaming ensures that all fields that are ever defined have an lowerCamelCase name.
func (checker *NamingChecker) CheckFieldNaming(method *ast.AssignStatement) {
	assignedField := method.Target
	identifier, isIdentifier := assignedField.(*ast.Identifier)
	if !isIdentifier {
		return
	}
	if !isLowerCamelCase(identifier.Value) {
		checker.reportInvalidNode(identifier, MessageInvalidDeclarationName)
	}
}

// CheckMethodNamingAndImplicitParameters ensures that a methods name is lowerCamelCase and that
// its parameters have explicit names if their type occurs more than once in the ParameterList.
// Meaning that when the ParameterList contains two numbers: '(number, number x)', both of the
// parameters need an explicit name: '(number x, number y)'.
func (checker *NamingChecker) CheckMethodNamingAndImplicitParameters(method *ast.MethodDeclaration) {
	if !isLowerCamelCase(method.Name.Value) {
		checker.reportInvalidNode(method, MessageInvalidDeclarationName)
	}
	checker.ensureExplicitParameterNamingOnDuplicateTypes(method.Parameters)
}

func (checker *NamingChecker) ensureExplicitParameterNamingOnDuplicateTypes(parameters ast.ParameterList) {
	parameterTypeNames := map[string] bool{}
	for _, parameter := range parameters {
		if parameterTypeNames[parameter.Type.NonGenericName()] {
			checker.reportInvalidNode(parameter, MessageImplicitParameterName)
		}
		parameterTypeNames[parameter.Type.NonGenericName()] = true
	}
}

// CheckParameterNaming ensures that a parameter is named lowerCamelCase.
func (checker *NamingChecker) CheckParameterNaming(parameter *ast.Parameter) {
	if isLowerCamelCase(parameter.Name.Value) {
		return
	}
	if parameter.IsNamedAfterType() {
		return
	}
	checker.reportInvalidNode(parameter, MessageInvalidDeclarationName)
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
