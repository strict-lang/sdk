package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/syntaxtree"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

func (parsing *Parsing) parseMethodDeclaration() (*syntaxtree.MethodDeclaration, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.MethodKeyword); err != nil {
		return nil, err
	}
	declaration, err := parsing.parseMethodSignature()
	if err != nil {
		return nil, err
	}
	var body syntaxtree.Node
	if token.OperatorValue(parsing.token()) == token.ArrowOperator {
		body, err = parsing.parseAssignedMethodExpression()
	} else {
		parsing.skipEndOfStatement()
		body, err = parsing.parseMethodBody(declaration.methodName.Value)
	}
	if err != nil {
		return nil, err
	}
	return &syntaxtree.MethodDeclaration{
		Type:         declaration.returnTypeName,
		Name:         declaration.methodName,
		Body:         body,
		Parameters:   declaration.parameters,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

type methodDeclaration struct {
	returnTypeName syntaxtree.TypeName
	methodName     *syntaxtree.Identifier
	parameters     syntaxtree.ParameterList
}

func (parsing *Parsing) parseMethodBody(methodName string) (node syntaxtree.Node, err error) {
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

func (parsing *Parsing) parseOptionalReturnTypeName() (syntaxtree.TypeName, error) {
	if parsing.isLookingAtOperator(token.LeftParenOperator) {
		return &syntaxtree.ConcreteTypeName{
			Name:         "void",
			NodePosition: parsing.createPosition(parsing.offset()),
		}, nil
	}
	return parsing.parseTypeName()
}

func (parsing *Parsing) parseAssignedMethodExpression() (syntaxtree.Node, error) {
	if err := parsing.skipOperator(token.ArrowOperator); err != nil {
		return nil, err
	}
	beginPosition := parsing.offset()
	statement := parsing.parseStatement()
	if expression, isExpression := statement.(*syntaxtree.ExpressionStatement); isExpression {
		return &syntaxtree.ReturnStatement{
			NodePosition: parsing.createPosition(beginPosition),
			Value:        expression,
		}, nil
	}
	return statement, nil
}

func (parsing *Parsing) parseParameterList() (parameters syntaxtree.ParameterList, err error) {
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

func (parsing *Parsing) parseParameter() (*syntaxtree.Parameter, error) {
	beginOffset := parsing.offset()
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	if next := parsing.token(); token.IsIdentifierToken(next) {
		idNameBegin := parsing.offset()
		parsing.advance()
		return &syntaxtree.Parameter{
			Name: &syntaxtree.Identifier{
				Value:        next.Value(),
				NodePosition: parsing.createPosition(idNameBegin),
			},
			Type:         typeName,
			NodePosition: parsing.createPosition(beginOffset),
		}, nil
	}
	return parsing.createTypeNamedParameter(beginOffset, typeName), nil
}

func (parsing *Parsing) createTypeNamedParameter(beginOffset source.Offset, typeName syntaxtree.TypeName) *syntaxtree.Parameter {
	return &syntaxtree.Parameter{
		Type: typeName,
		Name: &syntaxtree.Identifier{
			Value: typeName.NonGenericName(),
		},
		NodePosition: parsing.createPosition(beginOffset),
	}
}
