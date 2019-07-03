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
			Token:    nameToken,
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
		Name:      ast.NewIdentifier(nameToken.Value()),
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
				Token:    nextToken,
				Expected: "',' or ')'",
			}
		}
		operator := nextToken.(*token.OperatorToken).Operator
		if operator == token.RightParenOperator {
			break
		}
		if operator != token.CommaOperator {
			return arguments, &UnexpectedTokenError{
				Token:    nextToken,
				Expected: "',' or ')'",
			}
		}
	}
	return arguments, nil
}

func (parser *Parser) ParseIfStatement() ast.Node {
	return nil
}

func (parser *Parser) ParseForStatement() ast.Node {
	return nil
}

func (parser *Parser) ParseYieldStatement() ast.Node {
	return nil
}

func (parser *Parser) ParseReturnStatement() ast.Node {
	return nil
}

var keywordStatementParsers = map[token.Keyword]func(*Parser) ast.Node{
}

func (parser *Parser) keywordStatementParser(keyword token.Keyword) (func() ast.Node, bool) {
	switch keyword {
	case token.IfKeyword:
		return parser.ParseIfStatement, true
	case token.ForKeyword:
		return parser.ParseForStatement, true
	case token.YieldKeyword:
		return parser.ParseYieldStatement, true
	case token.ReturnKeyword:
		return parser.ParseReturnStatement, true
	}
	return nil, false
}

func (parser *Parser) ParseKeywordStatement(keyword token.Keyword) ast.Node {
	function, ok := parser.keywordStatementParser(keyword)
	if ok {
		return function()
	}
	parser.reportError(&UnexpectedTokenError{
		Token:    parser.tokens.Peek(),
		Expected: "statement begin",
	})
	return &ast.InvalidStatement{}
}

func (parser *Parser) reportInvalidStatement() {

}

func (parser *Parser) parseAssignStatement(leftHandSide ast.Node) (ast.Node, error) {
	return nil, nil
}

func (parser *Parser) ParseInstructionStatement() (ast.Node, error) {
	leftHandSide, err := parser.ParseLeftHandSide()
	if err != nil {
		return nil, err
	}
	nextToken := parser.tokens.Pull()
	switch operator := token.OperatorValue(nextToken); {
	case operator.IsAssign():
		return parser.parseAssignStatement(leftHandSide)
	}
	return nil, nil
}

func (parser *Parser) ParseStatement() ast.Node {
	switch peek := parser.tokens.Peek(); {
	case token.IsKeywordToken(peek):
		return parser.ParseKeywordStatement(peek.(*token.KeywordToken).Keyword)
	case token.IsIdentifierToken(peek),
		token.IsOperatorToken(peek),
		token.IsLiteralToken(peek):

		return parser.ParseInstructionStatement()
	default:
		parser.reportInvalidStatement()
		return &ast.InvalidStatement{}
	}
}

func (parser *Parser) ParseStatementSequence() []ast.Node {
	var statements []ast.Node
	for {
		next := parser.tokens.Pull()
		if next.Indent() != parser.block.Indent {
			break
		}
		statements = append(statements, parser.ParseStatement())
	}
	return statements
}
