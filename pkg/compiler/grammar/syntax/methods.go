package syntax

import (
	"strict.dev/sdk/pkg/compiler/diagnostic"
	"strict.dev/sdk/pkg/compiler/grammar/token"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/input"
)

const voidReturnType = `void`

type parsedMethod struct {
	name   string
	isVoid bool
}

func (parsing *Parsing) parseMethodDeclaration() *tree.MethodDeclaration {
	parsing.beginStructure(tree.MethodDeclarationNodeKind)
	parsing.skipKeyword(token.MethodKeyword)
	signature := parsing.parseMethodSignature()
	parsing.updateCurrentMethod(signature)
	return &tree.MethodDeclaration{
		Type:       signature.returnTypeName,
		Name:       signature.name,
		Parameters: signature.parameters,
		Body:       parsing.parseMethodBody(),
		Region:     parsing.completeStructure(tree.MethodDeclarationNodeKind),
	}
}

func (parsing *Parsing) updateCurrentMethod(signature methodSignature) {
	parsing.currentMethod = parsedMethod{
		name:   signature.name.Value,
		isVoid: isVoidType(signature.returnTypeName),
	}
}

func isVoidType(name tree.TypeName) bool {
	if concrete, isConcrete := name.(*tree.ConcreteTypeName); isConcrete {
		return concrete.Name == voidReturnType
	}
	return false
}

type methodSignature struct {
	name           *tree.Identifier
	parameters     tree.ParameterList
	returnTypeName tree.TypeName
}

func (parsing *Parsing) parseMethodBody() tree.Node {
	if token.OperatorValue(parsing.token()) == token.ArrowOperator {
		return parsing.parseAssignedMethodBody()
	}
	parsing.skipEndOfStatement()
	return parsing.parseMethodBlockBody()
}

func (parsing *Parsing) parseMethodBlockBody() tree.Node {
	return parsing.parseStatementBlock()
}

func (parsing *Parsing) parseMethodSignature() methodSignature {
	return methodSignature{
		returnTypeName: parsing.parseOptionalReturnTypeName(),
		name:           parsing.parseIdentifier(),
		parameters:     parsing.parseParameterListWithParens(),
	}
}

func (parsing *Parsing) parseOptionalReturnTypeName() tree.TypeName {
	if parsing.isLookingAtOperator(token.LeftParenOperator) {
		return &tree.ConcreteTypeName{
			Name:   "Void",
			Region: input.CreateRegion(parsing.offset(), parsing.offset()),
		}
	}
	return parsing.parseTypeName()
}

func (parsing *Parsing) parseAssignedMethodBody() tree.Node {
	parsing.skipOperator(token.ArrowOperator)
	statement := parsing.parseStatement()
	return &tree.StatementBlock{
		Children: []tree.Statement{
			parsing.replaceNodeWithReturnIfExpression(statement),
		},
		Region:   statement.Locate(),
	}
}

func (parsing *Parsing) replaceNodeWithReturnIfExpression(node tree.Node) tree.Node {
	if parsing.isReturningVoid() {
		return node
	}
	if expression, isExpression := node.(*tree.ExpressionStatement); isExpression {
		return &tree.ReturnStatement{
			Region: node.Locate(),
			Value:  expression.Expression,
		}
	}
	return node
}

func (parsing *Parsing) isReturningVoid() bool {
	return parsing.currentMethod.isVoid
}

func (parsing *Parsing) parseParameterListWithParens() tree.ParameterList {
	parsing.skipOperator(token.LeftParenOperator)
	defer parsing.skipOperator(token.RightParenOperator)
	return parsing.parseParameterList()
}

func (parsing *Parsing) parseParameterList() (parameters tree.ParameterList) {
	for !parsing.isAtEndOfParameterList() {
		parameters = append(parameters, parsing.parseParameter())
		parsing.consumeTokenAfterParameter()
	}
	return parameters
}

func (parsing *Parsing) consumeTokenAfterParameter() {
	if token.HasOperatorValue(parsing.token(), token.CommaOperator) {
		parsing.advance()
		return
	}
	parsing.expectEndOfParameterList()
}

func (parsing *Parsing) expectEndOfParameterList() {
	if !parsing.isAtEndOfParameterList() {
		parsing.throwError(&diagnostic.RichError{
			Error: &diagnostic.UnexpectedTokenError{
				Expected: ")",
				Received: parsing.token().Value(),
			},
			CommonReasons: []string{
				"The ParameterList is left open",
				"A Parameter declaration is invalid",
			},
		})
	}
}

func (parsing *Parsing) isAtEndOfParameterList() bool {
	operator := token.OperatorValue(parsing.token())
	return operator == token.RightParenOperator
}

func (parsing *Parsing) parseParameter() *tree.Parameter {
	parsing.beginStructure(tree.ParameterNodeKind)
	return &tree.Parameter{
		Type:   parsing.parseTypeName(),
		Name:   parsing.parseParameterName(),
		Region: parsing.completeStructure(tree.ParameterNodeKind),
	}
}

func (parsing *Parsing) parseParameterName() *tree.Identifier {
	currentToken := parsing.token()
	if token.IsIdentifierToken(currentToken) {
		return parsing.parseIdentifier()
	}
	parsing.throwError(newMissingParameterNameError())
	return nil
}

func newMissingParameterNameError() *diagnostic.RichError {
	return &diagnostic.RichError{
		Error:         &diagnostic.SpecificError{
			Message: "Name of the parameter is missing",
		},
		CommonReasons: []string{
			"The parameters type was not specified prior to the name",
		},
	}
}


