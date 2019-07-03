package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

// ParseMethodDeclaration parses a method declaration.
func (parser *Parser) ParseMethodDeclaration() (*ast.Method, error) {
	if last := parser.tokens.Last(); !token.HasKeywordValue(last, token.MethodKeyword) {
		return nil, &UnexpectedTokenError{
			Token:    last,
			Expected: token.MethodKeyword.String(),
		}
	}
	parser.tokens.Pull()
	returnTypeName, err := parser.ParseTypeName()
	if err != nil {
		return nil, err
	}
	methodName, err := parser.expectAnyIdentifier()
	if err != nil {
		return nil, err
	}
	parameters, err := parser.parseParameterList()
	if err != nil {
		return nil, err
	}
	body := parser.ParseStatementBlock()
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
		if token.OperatorValue(parser.tokens.Pull()) == token.RightParenOperator {
			break
		}
		parameter, err := parser.parseParameter()
		if err != nil {
			return parameters, err
		}
		parameters = append(parameters, parameter)
		switch next := parser.tokens.Pull(); {
		case token.OperatorValue(next) == token.CommaOperator:
			continue
		case token.OperatorValue(next) != token.RightParenOperator:
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
	if next := parser.tokens.Peek(); token.IsIdentifierToken(next) {
		parser.tokens.Pull()
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
