package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

func (parsing *Parsing) parseMethodDeclaration() (*ast.MethodDeclaration, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.MethodKeyword); err != nil {
		return nil, err
	}
	declaration, err := parsing.parseMethodSignature()
	if err != nil {
		return nil, err
	}
	var body ast.Node
	if token.OperatorValue(parsing.token()) == token.ArrowOperator {
		body, err = parsing.parseMethodAssignment()
	} else {
		parsing.skipEndOfStatement()
		body, err = parsing.parseMethodBody(declaration.methodName.Value)
	}
	if err != nil {
		return nil, err
	}
	return &ast.MethodDeclaration{
		Type:         declaration.returnTypeName,
		Name:         declaration.methodName,
		Body:         body,
		Parameters:   declaration.parameters,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

type methodDeclaration struct {
	returnTypeName ast.TypeName
	methodName     *ast.Identifier
	parameters     ast.ParameterList
}

func (parsing *Parsing) parseMethodBody(methodName string) (node ast.Node, err error) {
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

func (parsing *Parsing) parseOptionalReturnTypeName() (ast.TypeName, error) {
	if parsing.isLookingAtOperator(token.LeftParenOperator) {
		return &ast.ConcreteTypeName{
			Name:         "void",
			NodePosition: parsing.createPosition(parsing.offset()),
		}, nil
	}
	return parsing.parseTypeName()
}

func (parsing *Parsing) parseMethodAssignment() (ast.Node, error) {
	if err := parsing.skipOperator(token.ArrowOperator); err != nil {
		return nil, err
	}
	beginPosition := parsing.offset()
	statement := parsing.parseStatement()
	if expression, isExpression := statement.(*ast.ExpressionStatement); isExpression {
		return &ast.ReturnStatement{
			NodePosition: parsing.createPosition(beginPosition),
			Value:        expression,
		}, nil
	}
	return statement, nil
}

func (parsing *Parsing) parseParameterList() (parameters ast.ParameterList, err error) {
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

func (parsing *Parsing) parseParameter() (*ast.Parameter, error) {
	beginOffset := parsing.offset()
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	if next := parsing.token(); token.IsIdentifierToken(next) {
		idNameBegin := parsing.offset()
		parsing.advance()
		return &ast.Parameter{
			Name: &ast.Identifier{
				Value:        next.Value(),
				NodePosition: parsing.createPosition(idNameBegin),
			},
			Type:         typeName,
			NodePosition: parsing.createPosition(beginOffset),
		}, nil
	}
	return parsing.createTypeNamedParameter(beginOffset, typeName), nil
}

func (parsing *Parsing) createTypeNamedParameter(beginOffset source.Offset, typeName ast.TypeName) *ast.Parameter {
	return &ast.Parameter{
		Type: typeName,
		Name: &ast.Identifier{
			Value: typeName.NonGenericName(),
		},
		NodePosition: parsing.createPosition(beginOffset),
	}
}
