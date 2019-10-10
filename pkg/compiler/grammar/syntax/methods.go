package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
	"gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

func (parsing *Parsing) parseMethodDeclaration() (*tree.MethodDeclaration, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.MethodKeyword); err != nil {
		return nil, err
	}
	declaration, err := parsing.parseMethodSignature()
	if err != nil {
		return nil, err
	}
	var body tree.Node
	if token.OperatorValue(parsing.token()) == token.ArrowOperator {
		body, err = parsing.parseAssignedMethodExpression()
	} else {
		parsing.skipEndOfStatement()
		body, err = parsing.parseMethodBody(declaration.methodName.Value)
	}
	if err != nil {
		return nil, err
	}
	return &tree.MethodDeclaration{
		Type:         declaration.returnTypeName,
		Name:         declaration.methodName,
		Body:         body,
		Parameters:   declaration.parameters,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

type methodDeclaration struct {
	returnTypeName tree.TypeName
	methodName     *tree.Identifier
	parameters     tree.ParameterList
}

func (parsing *Parsing) parseMethodBody(methodName string) (node tree.Node, err error) {
	parsing.currentMethodName = methodName
	node, err = parsing.parseStatementBlock()
	parsing.currentMethodName = notParsingMethod
	return
}

func (parsing *Parsing) parseMethodSignature() (declaration methodDeclaration,
	err error) {
	declaration.returnTypeName, err = parsing.parseOptionalReturnTypeName()
	if err != nil {
		return methodDeclaration{}, err
	}
	declaration.methodName, err = parsing.expectAnyIdentifier()
	if err != nil {
		return methodDeclaration{}, err
	}
	parsing.advance()
	declaration.parameters, err = parsing.parseParameterList()
	if err != nil {
		return methodDeclaration{}, err
	}
	return
}

func (parsing *Parsing) parseOptionalReturnTypeName() (tree.TypeName, error) {
	if parsing.isLookingAtOperator(token.LeftParenOperator) {
		return &tree.ConcreteTypeName{
			Name:         "void",
			NodePosition: parsing.createPosition(parsing.offset()),
		}, nil
	}
	return parsing.parseTypeName()
}

func (parsing *Parsing) parseAssignedMethodExpression() (tree.Node, error) {
	if err := parsing.skipOperator(token.ArrowOperator); err != nil {
		return nil, err
	}
	beginPosition := parsing.offset()
	statement := parsing.parseStatement()
	if expression, isExpression := statement.(*tree.ExpressionStatement); isExpression {
		return &tree.ReturnStatement{
			NodePosition: parsing.createPosition(beginPosition),
			Value:        expression,
		}, nil
	}
	return statement, nil
}

func (parsing *Parsing) parseParameterList() (parameters tree.ParameterList, err error) {
	if err := parsing.skipOperator(token.LeftParenOperator); err != nil {
		return nil, err
	}
	for {
		if token.OperatorValue(parsing.token()) == token.RightParenOperator {
			parsing.advance()
			break
		}
		parameter, err := parsing.parseParameter()
		if err != nil {
			return parameters, err
		}
		parameters = append(parameters, parameter)
		switch next := parsing.token(); {
		case token.OperatorValue(next) == token.CommaOperator:
			parsing.advance()
			continue
		case token.OperatorValue(next) != token.RightParenOperator:
			parsing.advance()
			return parameters, &UnexpectedTokenError{
				Token:    next,
				Expected: "end of method parameter list",
			}
		}
	}
	return parameters, nil
}

func (parsing *Parsing) parseParameter() (*tree.Parameter, error) {
	beginOffset := parsing.offset()
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	if next := parsing.token(); token.IsIdentifierToken(next) {
		idNameBegin := parsing.offset()
		parsing.advance()
		return &tree.Parameter{
			Name: &tree.Identifier{
				Value:        next.Value(),
				NodePosition: parsing.createPosition(idNameBegin),
			},
			Type:         typeName,
			NodePosition: parsing.createPosition(beginOffset),
		}, nil
	}
	return parsing.createTypeNamedParameter(beginOffset, typeName), nil
}

func (parsing *Parsing) createTypeNamedParameter(beginOffset input.Offset, typeName tree.TypeName) *tree.Parameter {
	return &tree.Parameter{
		Type: typeName,
		Name: &tree.Identifier{
			Value: typeName.NonGenericName(),
		},
		NodePosition: parsing.createPosition(beginOffset),
	}
}
