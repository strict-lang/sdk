package parser

import (
	"github.com/BenjaminNitschke/Strict/pkg/ast"
)

// ParseConditionalStatement parses a conditional statement and its optional else-clause.
// Conditional statements are also called 'if'-statements. They have an expression, a body
// and an optional else clause which itself has a body. If the body-node of an else-clause
// is another conditional statement, it is referred to as an 'else-if'.
func (parser *Parser) ParseConditionalStatement() (ast.ConditionalStatement, error) {
	if ok, err := parser.skipKeyword(token.IfKeyword); !ok {
		return ast.ConditionalStatement{}, err
	}
	condition, err := parser.ParseExpression()
	if err != nil {
		return ast.ConditionalStatement{}, err
	}
	if ok, err := parser.expectOperator(token.ColonOperator); !ok {
		return ast.ConditionalStatement{}, err
	}
	body, err := parser.ParseStatements()
	if err != nil {
		return ast.ConditionalStatement{}, err
	}
	if !parser.isLookingAtKeyword(token.ElseKeyword) {
		return ast.ConditionalStatement{
			body:      body,
			condition: condition,
		}, nil
	}

	// Remove the else clause from the token queue
	parser.tokens.Pull()
	elseClause, err := parser.ParseElseClause()
	if err != nil {
		return ast.ConditionalStatement{}, err
	}
	return ast.ConditionalStatement{
		body:       body,
		condition:  condition,
		elseClause: elseClause,
	}
}

// ParseElseClause parses the else-clause of an conditional statement.
func (parser *Parser) ParseElseClause() (ast.Node, error) {
	if parser.isLookingAtKeyword(token.IfKeyword) {
		return parser.ParseConditionalStatement()
	}
	return parser.ParseStatements()
}
