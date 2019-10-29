package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (parsing *Parsing) parseMethodDeclaration() *tree.MethodDeclaration {
	beginOffset := parsing.offset()
	parsing.skipKeyword(token.MethodKeyword)
	signature := parsing.parseMethodSignature()
	parsing.currentMethodName = signature.name.Value
	return &tree.MethodDeclaration{
		Type:       signature.returnTypeName,
		Name:       signature.name,
		Parameters: signature.parameters,
		Body:       parsing.parseMethodBody(),
		Region:     parsing.createRegion(beginOffset),
	}
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
			Name:   "void",
			Region: parsing.createRegion(parsing.offset()),
		}
	}
	return parsing.parseTypeName()
}

func (parsing *Parsing) parseAssignedMethodBody() tree.Node {
	parsing.skipOperator(token.ArrowOperator)
	statement := parsing.parseStatement()
	return replaceNodeWithReturnIfExpression(statement)
}

func replaceNodeWithReturnIfExpression(node tree.Node) tree.Node {
	if expression, isExpression := node.(*tree.ExpressionStatement); isExpression {
		return &tree.ReturnStatement{
			Region: node.Locate(),
			Value:  expression,
		}
	}
	return node
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
	next := parsing.pullToken()
	if token.HasOperatorValue(next, token.CommaOperator) {
		return
	}
	parsing.expectEndOfParameterList()
}

func (parsing *Parsing) expectEndOfParameterList() {
	if !parsing.isAtEndOfParameterList() {
		parsing.throwError(&UnexpectedTokenError{
			Token:    parsing.token(),
			Expected: "end of method parameter list",
		})
	}
}

func (parsing *Parsing) isAtEndOfParameterList() bool {
	operator := token.OperatorValue(parsing.token())
	return operator != token.RightParenOperator
}

func (parsing *Parsing) parseParameter() *tree.Parameter {
	beginOffset := parsing.offset()
	return &tree.Parameter{
		Type:   parsing.parseTypeName(),
		Name:   parsing.parseIdentifier(),
		Region: parsing.createRegion(beginOffset),
	}
}
