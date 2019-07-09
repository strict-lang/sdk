package parser

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

// ParseMethodDeclaration parses a method declaration.
func (parser *Parser) ParseMethodDeclaration() (*ast.Method, error) {
	if err := parser.skipKeyword(token.MethodKeyword); err != nil {
		return nil, err
	}
	returnTypeName, err := parser.ParseTypeName()
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
	body, err := parser.ParseStatementBlock()
	if err != nil {
		return nil, err
	}
	return &ast.Method{
		Type:       returnTypeName,
		Name:       methodName,
		Parameters: parameters,
		Body:       body,
	}, nil
}

func (parser *Parser) parseParameterList() ([]ast.Parameter, error) {
	if err := parser.skipOperator(token.LeftParenOperator); err != nil {
		return nil, err
	}
	var parameters []ast.Parameter
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
				Expected: ", or )",
			}
		}
	}
	return parameters, nil
}

func (parser *Parser) parseParameter() (ast.Parameter, error) {
	typeName, err := parser.ParseTypeName()
	if err != nil {
		return ast.Parameter{}, err
	}
	if next := parser.token(); token.IsIdentifierToken(next) {
		parser.advance()
		return ast.Parameter{
			Name: ast.Identifier{Value: next.Value()},
			Type: typeName,
		}, nil
	}
	return typeNamedParameter(typeName), nil
}

func typeNamedParameter(typeName ast.TypeName) ast.Parameter {
	return ast.Parameter{
		Type: typeName,
		Name: ast.Identifier{
			Value: typeName.NonGenericName(),
		},
	}
}
