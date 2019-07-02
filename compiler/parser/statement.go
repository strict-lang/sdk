package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

func (parser *Parser) ParseStatements() (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) ParseMethodCall() (ast.MethodCall, error) {
	nameToken := parser.tokens.Last()
	if !token.IsIdentifierToken(nameToken) {
		return ast.MethodCall{}, &UnexpectedTokenError{
			Token: nameToken,
			Expected: "an identifier",
		}
	}
 	if err := parser.skipOperator(token.LeftParenOperator); err != nil {
 		return ast.MethodCall{}, err
	}
 	arguments, err := parser.parseArgumentList()
 	if err != nil {
 		return ast.MethodCall{}, err
	}
 	return ast.MethodCall{
 		Arguments: arguments,
		Name: ast.NewIdentifier(nameToken.Value()),
	}, nil
}

func (parser *Parser) parseArgumentList() ([]ast.Node, error) {
	var arguments []ast.Node
	for {
		next, err := parser.ParseExpression()
		if err != nil {
			return arguments, err
		}
		arguments = append(arguments, next)
		nextToken := parser.tokens.Pull()
		if !token.IsOperatorToken(nextToken) {
			return arguments, &UnexpectedTokenError{
				Token: nextToken,
				Expected: "',' or ')'",
			}
		}
		operator := nextToken.(*token.OperatorToken).Operator
		if operator == token.RightParenOperator {
			break
		}
		if operator != token.CommaOperator {
			return arguments, &UnexpectedTokenError{
				Token: nextToken,
				Expected: "',' or ')'",
			}
		}
	}
	return arguments, nil
}
