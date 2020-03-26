package syntax

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/token"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

func (parsing *Parsing) parseImportStatementList() (imports []*tree.ImportStatement) {
	for token.HasKeywordValue(parsing.token(), token.ImportKeyword) {
		imports = append(imports, parsing.parseImportStatement())
	}
	return imports
}

func (parsing *Parsing) parseClassDeclaration() *tree.ClassDeclaration {
	parsing.beginStructure(tree.ClassDeclarationNodeKind)
	nodes := parsing.parseTopLevelNodes()
	return &tree.ClassDeclaration{
		Name:       parseFileName(parsing.unitName),
		Parameters: []*tree.ClassParameter{},
		SuperTypes: []tree.TypeName{},
		Children:   nodes,
		Region:     parsing.completeStructure(tree.ClassDeclarationNodeKind),
	}
}

func (parsing *Parsing) parseTestStatement() tree.Node {
	parsing.beginStructure(tree.TestStatementNodeKind)
	parsing.skipKeyword(token.TestKeyword)
	parsing.skipEndOfStatement()
	statements := parsing.parseStatementBlock()
	return &tree.TestStatement{
		Region:     parsing.completeStructure(tree.TestStatementNodeKind),
		MethodName: parsing.currentMethod.name,
		Body:       statements,
	}
}

func (parsing *Parsing) parseFieldDeclaration() tree.Node {
	parsing.beginStructure(tree.FieldDeclarationNodeKind)
	parsing.skipKeyword(token.HasKeyword)
	name := parsing.parseIdentifier()
	fieldType := parsing.parseTypeName()
	return &tree.FieldDeclaration{
		Region: parsing.completeStructure(tree.FieldDeclarationNodeKind),
		Name: name,
		TypeName: fieldType,
	}
}

func (parsing *Parsing) parseLetBinding() *tree.LetBinding {
	parsing.beginStructure(tree.LetBindingNodeKind)
	parsing.skipKeyword(token.LetKeyword)
	name := parsing.parseIdentifier()
	parsing.skipOperator(token.AssignOperator)
	value := parsing.parseExpression()
	parsing.skipEndOfStatement()
	return &tree.LetBinding{
		Region:     parsing.completeStructure(tree.LetBindingNodeKind),
		Name:       name,
		Expression: value,
	}
}

func (parsing *Parsing) parseLetBindingStatement() tree.Statement {
	parsing.beginStructure(tree.ExpressionStatementNodeKind)
	binding := parsing.parseLetBinding()
	_ = parsing.completeStructure(tree.ExpressionStatementNodeKind)
	return &tree.ExpressionStatement{
		Expression: binding,
	}
}

func (parsing *Parsing) parseImplementStatement() tree.Node {
	parsing.beginStructure(tree.ImplementStatementNodeKind)
	parsing.skipKeyword(token.ImplementKeyword)
	trait := parsing.parseTypeName()
	parsing.skipEndOfStatement()
	return &tree.ImplementStatement{
		Region: parsing.completeStructure(tree.ImplementStatementNodeKind),
		Trait:  trait,
	}
}

func (parsing *Parsing) parseGenericStatement() tree.Node {
	parsing.beginStructure(tree.GenericStatementNodeKind)
	parsing.skipKeyword(token.GenericKeyword)
	name := parsing.parseIdentifier()
	constraints := parsing.parseGenericConstraints()
	parsing.skipEndOfStatement()
	return &tree.GenericStatement{
		Region:      parsing.completeStructure(tree.GenericStatementNodeKind),
		Name:        name,
		Constraints: constraints,
	}
}

func (parsing *Parsing) parseGenericConstraints() []tree.TypeName {
	if !parsing.isLookingAtKeyword(token.IsKeyword) {
		return []tree.TypeName{}
	}
	parsing.skipKeyword(token.IsKeyword)
	return parsing.parseTypeNameList()
}

func (parsing *Parsing) parseTypeNameList() (names []tree.TypeName) {
	names = append(names, parsing.parseTypeName())
	for parsing.isLookingAtOperator(token.CommaOperator) {
		parsing.skipOperator(token.CommaOperator)
		names = append(names, parsing.parseTypeName())
	}
	return names
}