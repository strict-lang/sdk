package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/source"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

// ParseMethodDeclaration parses a method declaration.
func (parser *Parser) ParseMethodDeclaration() (*ast.MethodDeclaration, error) {
	beginOffset := parser.offset()
	if err := parser.skipKeyword(token.MethodKeyword); err != nil {
		return nil, err
	}
	returnTypeName, err := parser.parseOptionalReturnTypeName()
	if err != nil {
		return nil, err
	}
	methodName, err := parser.expectAnyIdentifier()
	if err != nil {
		return nil, err
	}
	parser.advance()
	parameters, err := parser.parseParameterList()
	if err != nil {
		return nil, err
	}
	parser.skipEndOfStatement()
	parser.currentMethodName = methodName.Value
	body, err := parser.ParseStatementBlock()
	parser.currentMethodName = notParsingMethod
	if err != nil {
		return nil, err
	}
	return &ast.MethodDeclaration{
		Type:         returnTypeName,
		Name:         methodName,
		Body:         body,
		Parameters:   parameters,
		NodePosition: parser.createPosition(beginOffset),
	}, nil
}

func (parser *Parser) parseOptionalReturnTypeName() (ast.TypeName, error) {
	if parser.isLookingAtOperator(token.LeftParenOperator) {
		return &ast.ConcreteTypeName{
			Name:         "void",
			NodePosition: parser.createPosition(parser.offset()),
		}, nil
	}
	return parser.ParseTypeName()
}

func (parser *Parser) parseParameterList() (parameters ast.ParameterList, err error) {
	if err := parser.skipOperator(token.LeftParenOperator); err != nil {
		return nil, err
	}
	for {
		if token.OperatorValue(parser.token()) == token.RightParenOperator {
			parser.advance()
			break
		}
		parameter, err := parser.parseParameter()
		if err != nil {
			return parameters, err
		}
		parameters = append(parameters, parameter)
		switch next := parser.token(); {
		case token.OperatorValue(next) == token.CommaOperator:
			parser.advance()
			continue
		case token.OperatorValue(next) != token.RightParenOperator:
			parser.advance()
			return parameters, &UnexpectedTokenError{
				Token:    next,
				Expected: "end of method parameter list",
			}
		}
	}
	return parameters, nil
}

func (parser *Parser) parseParameter() (*ast.Parameter, error) {
	beginOffset := parser.offset()
	typeName, err := parser.ParseTypeName()
	if err != nil {
		return nil, err
	}
	if next := parser.token(); token.IsIdentifierToken(next) {
		idNameBegin := parser.offset()
		parser.advance()
		return &ast.Parameter{
			Name: &ast.Identifier{
				Value:        next.Value(),
				NodePosition: parser.createPosition(idNameBegin),
			},
			Type:         typeName,
			NodePosition: parser.createPosition(beginOffset),
		}, nil
	}
	return parser.createTypeNamedParameter(beginOffset, typeName), nil
}

func (parser *Parser) createTypeNamedParameter(beginOffset source.Offset, typeName ast.TypeName) *ast.Parameter {
	return &ast.Parameter{
		Type: typeName,
		Name: &ast.Identifier{
			Value: typeName.NonGenericName(),
		},
		NodePosition: parser.createPosition(beginOffset),
	}
}
