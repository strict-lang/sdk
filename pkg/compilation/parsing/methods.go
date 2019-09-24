package parsing

import (
	source2 "gitlab.com/strict-lang/sdk/pkg/compilation/source"
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	token2 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

func (parsing *Parsing) parseMethodDeclaration() (*syntaxtree2.MethodDeclaration, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token2.MethodKeyword); err != nil {
		return nil, err
	}
	declaration, err := parsing.parseMethodSignature()
	if err != nil {
		return nil, err
	}
	var body syntaxtree2.Node
	if token2.OperatorValue(parsing.token()) == token2.ArrowOperator {
		body, err = parsing.parseAssignedMethodExpression()
	} else {
		parsing.skipEndOfStatement()
		body, err = parsing.parseMethodBody(declaration.methodName.Value)
	}
	if err != nil {
		return nil, err
	}
	return &syntaxtree2.MethodDeclaration{
		Type:         declaration.returnTypeName,
		Name:         declaration.methodName,
		Body:         body,
		Parameters:   declaration.parameters,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

type methodDeclaration struct {
	returnTypeName syntaxtree2.TypeName
	methodName     *syntaxtree2.Identifier
	parameters     syntaxtree2.ParameterList
}

func (parsing *Parsing) parseMethodBody(methodName string) (node syntaxtree2.Node, err error) {
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

func (parsing *Parsing) parseOptionalReturnTypeName() (syntaxtree2.TypeName, error) {
	if parsing.isLookingAtOperator(token2.LeftParenOperator) {
		return &syntaxtree2.ConcreteTypeName{
			Name:         "void",
			NodePosition: parsing.createPosition(parsing.offset()),
		}, nil
	}
	return parsing.parseTypeName()
}

func (parsing *Parsing) parseAssignedMethodExpression() (syntaxtree2.Node, error) {
	if err := parsing.skipOperator(token2.ArrowOperator); err != nil {
		return nil, err
	}
	beginPosition := parsing.offset()
	statement := parsing.parseStatement()
	if expression, isExpression := statement.(*syntaxtree2.ExpressionStatement); isExpression {
		return &syntaxtree2.ReturnStatement{
			NodePosition: parsing.createPosition(beginPosition),
			Value:        expression,
		}, nil
	}
	return statement, nil
}

func (parsing *Parsing) parseParameterList() (parameters syntaxtree2.ParameterList, err error) {
	if err := parsing.skipOperator(token2.LeftParenOperator); err != nil {
		return nil, err
	}
	for {
		if token2.OperatorValue(parsing.token()) == token2.RightParenOperator {
			parsing.advance()
			break
		}
		parameter, err := parsing.parseParameter()
		if err != nil {
			return parameters, err
		}
		parameters = append(parameters, parameter)
		switch next := parsing.token(); {
		case token2.OperatorValue(next) == token2.CommaOperator:
			parsing.advance()
			continue
		case token2.OperatorValue(next) != token2.RightParenOperator:
			parsing.advance()
			return parameters, &UnexpectedTokenError{
				Token:    next,
				Expected: "end of method parameter list",
			}
		}
	}
	return parameters, nil
}

func (parsing *Parsing) parseParameter() (*syntaxtree2.Parameter, error) {
	beginOffset := parsing.offset()
	typeName, err := parsing.parseTypeName()
	if err != nil {
		return nil, err
	}
	if next := parsing.token(); token2.IsIdentifierToken(next) {
		idNameBegin := parsing.offset()
		parsing.advance()
		return &syntaxtree2.Parameter{
			Name: &syntaxtree2.Identifier{
				Value:        next.Value(),
				NodePosition: parsing.createPosition(idNameBegin),
			},
			Type:         typeName,
			NodePosition: parsing.createPosition(beginOffset),
		}, nil
	}
	return parsing.createTypeNamedParameter(beginOffset, typeName), nil
}

func (parsing *Parsing) createTypeNamedParameter(beginOffset source2.Offset, typeName syntaxtree2.TypeName) *syntaxtree2.Parameter {
	return &syntaxtree2.Parameter{
		Type: typeName,
		Name: &syntaxtree2.Identifier{
			Value: typeName.NonGenericName(),
		},
		NodePosition: parsing.createPosition(beginOffset),
	}
}
