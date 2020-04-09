package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (parsing *Parsing) parseImportStatement() *tree.ImportStatement {
	parsing.beginStructure(tree.ImportStatementNodeKind)
	parsing.skipKeyword(token.ImportKeyword)
	switch next := parsing.token(); {
	case token.IsStringLiteralToken(next):
		return parsing.completeFileImport()
	case token.IsIdentifierToken(next):
		return parsing.completeIdentifierChainImport()
	default:
		parsing.reportInvalidImportTarget()
		return nil
	}
}

func (parsing *Parsing) reportInvalidImportTarget() {
	parsing.throwError(&diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "file or path to class",
			Received: parsing.token().Value(),
		},
		CommonReasons: []string{"The import target is invalid"},
	})
}

func (parsing *Parsing) completeIdentifierChainImport() *tree.ImportStatement {
	var chain []string
	for token.IsIdentifierToken(parsing.token()) {
		chain = append(chain, parsing.token().Value())
		parsing.advance()
		if token.HasOperatorValue(parsing.token(), token.DotOperator) {
			parsing.advance()
		}
	}
	parsing.skipEndOfStatement()
	return &tree.ImportStatement{
		Target: &tree.IdentifierChainImport{Chain: chain},
		Alias:  nil,
		Region: parsing.completeStructure(tree.ImportStatementNodeKind),
	}
}

func (parsing *Parsing) completeFileImport() *tree.ImportStatement {
	target := &tree.FileImport{Path: parsing.token().Value()}
	parsing.advance()
	if !token.HasKeywordValue(parsing.token(), token.AsKeyword) {
		parsing.skipEndOfStatement()
		return &tree.ImportStatement{
			Target: target,
			Region: parsing.completeStructure(tree.ImportStatementNodeKind),
		}
	}
	parsing.advance()
	return parsing.completeFileImportWithAlias(target)
}

func (parsing *Parsing) completeFileImportWithAlias(target tree.ImportTarget) *tree.ImportStatement {
	alias := parsing.parseImportAlias()
	parsing.skipEndOfStatement()
	return &tree.ImportStatement{
		Target: target,
		Alias:  alias,
		Region: parsing.completeStructure(tree.ImportStatementNodeKind),
	}
}

func (parsing *Parsing) parseImportAlias() *tree.Identifier {
	return parsing.parseIdentifier()
}

func (parsing *Parsing) parseOptionalAssignValue() tree.Expression {
	parsing.skipOperator(token.AssignOperator)
	return parsing.parseExpression()
}
