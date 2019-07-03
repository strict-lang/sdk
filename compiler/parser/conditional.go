package parser

import (
	"github.com/BenjaminNitschke/Strict/compiler/ast"
	"github.com/BenjaminNitschke/Strict/compiler/token"
)

// ParseConditionalStatement parses a conditional statement and its optional else-clause.
// Conditional statements are also called 'if'-statements. They have an expression, a body
// and an optional else clause which itself has a body. If the body-node of an else-clause
// is another conditional statement, it is referred to as an 'else-if'.
func (parser *Parser) ParseConditionalStatement() (*ast.ConditionalStatement, error) {
	if err := parser.skipKeyword(token.IfKeyword); err != nil {
		return &ast.ConditionalStatement{}, err
	}
	condition, err := parser.ParseExpression()
	if err != nil {
		return &ast.ConditionalStatement{}, err
	}
	if err := parser.expectOperator(token.ColonOperator); err != nil {
		return &ast.ConditionalStatement{}, err
	}
	body := parser.ParseStatementBlock()
	if !parser.isLookingAtKeyword(token.ElseKeyword) {
		return &ast.ConditionalStatement{
			Body:      body,
			Condition: condition,
		}, nil
	}

	// Remove the else clause from the token queue
	parser.advance()
	elseClause, err := parser.ParseElseClause()
	if err != nil {
		return &ast.ConditionalStatement{}, err
	}
	return &ast.ConditionalStatement{
		Else:      elseClause,
		Body:      body,
		Condition: condition,
	}, nil
}

// ParseElseClause parses the else-clause of an conditional statement.
func (parser *Parser) ParseElseClause() (ast.Node, error) {
	if parser.isLookingAtKeyword(token.IfKeyword) {
		return parser.ParseConditionalStatement()
	}
	return parser.ParseStatementBlock(), nil
}
