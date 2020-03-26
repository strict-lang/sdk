package syntax

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/diagnostic"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type keywordStatementParser func(*Parsing) tree.Node

var keywordStatementParserTable map[token.Keyword]keywordStatementParser
func init() {
	keywordStatementParserTable = map[token.Keyword]keywordStatementParser{
		token.IfKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseConditionalStatement()
		},
		token.ForKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseLoopStatement()
		},
		token.YieldKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseYieldStatement()
		},
		token.BreakKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseBreakStatement()
		},
		token.ReturnKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseReturnStatement()
		},
		token.ImportKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseImportStatement()
		},
		token.AssertKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseAssertStatement()
		},
		token.TestKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseTestStatement()
		},
		token.MethodKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseMethodDeclaration()
		},
		token.LetKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseLetBindingStatement()
		},
		token.ImplementKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseImplementStatement()
		},
		token.GenericKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseGenericStatement()
		},
		token.HasKeyword: func(parsing *Parsing) tree.Node {
			return parsing.parseFieldDeclaration()
		},
	}
}

func (parsing *Parsing) findKeywordStatementParser(
	keyword token.Keyword) (keywordStatementParser, bool) {
	parser, found := keywordStatementParserTable[keyword]
	return parser, found
}

func (parsing *Parsing) parseKeywordStatement(keyword token.Keyword) tree.Node {
	function, ok := parsing.findKeywordStatementParser(keyword)
	if ok {
		return function(parsing)
	}
	parsing.throwError(&diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "begin of statement",
			Received: parsing.token().Value(),
		},
	})
	return nil
}

func (parsing *Parsing) parseStatement() tree.Node {
	parsing.beginStructure(tree.UnknownNodeKind)
	defer parsing.completeStructure(tree.WildcardNodeKind) // Could have been modified

	switch current := parsing.token(); {
	case token.IsKeywordToken(current):
		return parsing.parseKeywordStatement(token.KeywordValue(current))
	case token.IsOperatorToken(current):
		fallthrough
	case token.IsIdentifierToken(current):
		fallthrough
	case token.IsLiteralToken(current):
		return parsing.parseInstructionStatement()
	default:
		parsing.throwError(newUnexpectedTokenError(current))
		return nil
	}
}

func newUnexpectedTokenError(token token.Token) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.UnexpectedTokenError{
			Expected: "begin of statement",
			Received: token.Value(),
		},
	}
}

func (parsing *Parsing) parseStatementSequence() (statements []tree.Statement) {
	for !token.IsEndOfFileToken(parsing.token()) {
		statement, shouldContinue := parsing.parseStatementInSequence()
		if !shouldContinue {
			break
		}
		if statement != nil {
			statements = append(statements, statement)
		}
	}
	return statements
}

func (parsing *Parsing) parseStatementInSequence() (statement tree.Statement, shouldContinue bool) {
	expectedIndent := parsing.block.Indent
	current := parsing.token()
	if token.IsEndOfStatementToken(current) {
		parsing.advance()
		return nil, true
	}
	if current.Indent() > expectedIndent {
		parsing.throwError(newInvalidIndentError(expectedIndent, current.Indent()))
		return nil, false
	}
	if current.Indent() < expectedIndent {
		return nil, false
	}
	statement = parsing.parseStatement()
	if _, ok := statement.(*tree.InvalidStatement); ok {
		return statement, false
	}
	return statement, true
}


func newInvalidIndentError(expected, received token.Indent) *diagnostic.RichError {
	return &diagnostic.RichError{
		Error: &diagnostic.InvalidIndentationError{
			Expected: fmt.Sprintf("%d", expected),
			Received: int(received),
		},
		CommonReasons: []string{
			"Tabs and spaces are mixed",
			"Code is not properly formatted",
		},
	}
}

// ParseStatementBlock parses a block of statements.
func (parsing *Parsing) parseStatementBlock() *tree.StatementBlock {
	parsing.beginStructure(tree.StatementBlockNodeKind)
	indent := parsing.token().Indent()
	if indent < parsing.block.Indent {
		parsing.throwError(newSmallerIndentError(indent))
	}
	parsing.openBlock(indent)
	statements := parsing.parseStatementSequence()
	parsing.closeBlock()
	return &tree.StatementBlock{
		Children: statements,
		Region:   parsing.completeStructure(tree.StatementBlockNodeKind),
	}
}
