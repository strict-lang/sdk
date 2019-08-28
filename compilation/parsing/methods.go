package parsing

import (
	"gitlab.com/strict-lang/sdk/compilation/ast"
	"gitlab.com/strict-lang/sdk/compilation/source"
	"gitlab.com/strict-lang/sdk/compilation/token"
)

// ParseMethodDeclaration parses a method declaration.
func (parsing *Parsing) ParseMethodDeclaration() (*ast.MethodDeclaration, error) {
	beginOffset := parsing.offset()
	if err := parsing.skipKeyword(token.MethodKeyword); err != nil {
		return nil, err
	}
	returnTypeName, err := parsing.parseOptionalReturnTypeName()
	if err != nil {
		return nil, err
	}
	methodName, err := parsing.expectAnyIdentifier()
	if err != nil {
		return nil, err
	}
	parsing.advance()
	parameters, err := parsing.parseParameterList()
	if err != nil {
		return nil, err
	}
	parsing.skipEndOfStatement()
	parsing.currentMethodName = methodName.Value
	body, err := parsing.ParseStatementBlock()
	parsing.currentMethodName = notParsingMethod
	if err != nil {
		return nil, err
	}
	return &ast.MethodDeclaration{
		Type:         returnTypeName,
		Name:         methodName,
		Body:         body,
		Parameters:   parameters,
		NodePosition: parsing.createPosition(beginOffset),
	}, nil
}

func (parsing *Parsing) parseOptionalReturnTypeName() (ast.TypeName, error) {
	if parsing.isLookingAtOperator(token.LeftParenOperator) {
		return &ast.ConcreteTypeName{
			Name:         "void",
			NodePosition: parsing.createPosition(parsing.offset()),
		}, nil
	}
	return parsing.ParseTypeName()
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
	typeName, err := parsing.ParseTypeName()
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
