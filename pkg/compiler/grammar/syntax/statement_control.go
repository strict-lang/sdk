package syntax

import (
	"github.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"github.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

// parseConditionalStatement parses a conditional statement and it's optional else-clause.
func (parsing *Parsing) parseConditionalStatement() *tree.ConditionalStatement {
	parsing.beginStructure(tree.ConditionalStatementNodeKind)
	parsing.skipKeyword(token.IfKeyword)
	condition := parsing.parseConditionalExpression()
	parsing.skipEndOfStatement()
	consequence := parsing.parseStatementBlock()
	return parsing.parseElseClauseIfPresent(condition, consequence)
}

func (parsing *Parsing) parseElseClauseIfPresent(
	condition tree.Expression,
	consequence *tree.StatementBlock) *tree.ConditionalStatement {

	if token.HasKeywordValue(parsing.token(), token.ElseKeyword) {
		return parsing.parseConditionalStatementWithAlternative(
			condition, consequence)
	}
	return &tree.ConditionalStatement{
		Condition:   condition,
		Consequence: consequence,
		Region:      parsing.completeStructure(tree.ConditionalStatementNodeKind),
	}
}

func (parsing *Parsing) parseConditionalStatementWithAlternative(
	condition tree.Expression,
	consequence *tree.StatementBlock) *tree.ConditionalStatement {

	parsing.advance()
	alternative := parsing.parseElseIfOrBlock()
	return &tree.ConditionalStatement{
		Condition:   condition,
		Consequence: consequence,
		Alternative: alternative,
		Region:      parsing.completeStructure(tree.ConditionalStatementNodeKind),
	}
}

func (parsing *Parsing) parseElseIfOrBlock() *tree.StatementBlock {
	if token.HasKeywordValue(parsing.token(), token.IfKeyword) {
		statement := parsing.parseConditionalStatement()
		return &tree.StatementBlock{
			Children: []tree.Statement{statement},
			Region:   statement.Region,
		}
	}
	parsing.skipEndOfStatement()
	return parsing.parseStatementBlock()
}

func (parsing *Parsing) parseLoopStatement() tree.Node {
	parsing.beginStructure(tree.ForEachLoopStatementNodeKind)
	parsing.skipKeyword(token.ForKeyword)
	field := parsing.parseIdentifier()
	parsing.skipKeyword(token.InKeyword)
	value := parsing.parseExpression()
	parsing.skipEndOfStatement()
	body := parsing.parseStatementBlock()
	return &tree.ForEachLoopStatement{
		Field:    field,
		Sequence: value,
		Body:     body,
		Region:   parsing.completeStructure(tree.ForEachLoopStatementNodeKind),
	}
}

func (parsing *Parsing) parseYieldStatement() *tree.YieldStatement {
	parsing.beginStructure(tree.YieldStatementNodeKind)
	parsing.skipKeyword(token.YieldKeyword)
	rightHandSide := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.YieldStatement{
		Value:  rightHandSide,
		Region: parsing.completeStructure(tree.YieldStatementNodeKind),
	}
}

func (parsing *Parsing) parseBreakStatement() *tree.BreakStatement {
	parsing.beginStructure(tree.BreakStatementNodeKind)
	parsing.skipKeyword(token.BreakKeyword)
	parsing.skipEndOfStatement()
	return &tree.BreakStatement{
		Region: parsing.completeStructure(tree.BreakStatementNodeKind),
	}
}

func (parsing *Parsing) parseAssertStatement() tree.Node {
	parsing.beginStructure(tree.AssertStatementNodeKind)
	parsing.skipKeyword(token.AssertKeyword)
	expression := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.AssertStatement{
		Region:     parsing.completeStructure(tree.AssertStatementNodeKind),
		Expression: expression,
	}
}

func (parsing *Parsing) parseReturnStatement() *tree.ReturnStatement {
	parsing.beginStructure(tree.ReturnStatementNodeKind)
	parsing.skipKeyword(token.ReturnKeyword)
	defer parsing.skipEndOfStatement()
	if token.IsEndOfStatementToken(parsing.token()) {
		parsing.advance()
		return &tree.ReturnStatement{
			Region: parsing.completeStructure(tree.ReturnStatementNodeKind),
		}
	}
	rightHandSide := parsing.parseExpression()
	return &tree.ReturnStatement{
		Value:  rightHandSide,
		Region: parsing.completeStructure(tree.ReturnStatementNodeKind),
	}
}
