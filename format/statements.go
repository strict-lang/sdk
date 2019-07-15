package format

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
	"gitlab.com/strict-lang/sdk/compiler/token"
)

func (printer *PrettyPrinter) printBlockStatement(block *ast.BlockStatement) {
	for _, statement := range block.Children {
		printer.appendIndent()
		printer.printNode(statement)
		printer.appendLineBreak()
	}
}

func (printer *PrettyPrinter) printExpressionStatement(statement *ast.ExpressionStatement) {
	printer.printNode(statement.Expression)
}

func (printer *PrettyPrinter) printConditionalStatement(statement *ast.ConditionalStatement) {
	printer.sawReturn = false
	printer.printIfHeader(statement)
	printer.indent.Open()
	printer.printNode(statement.Body)
	printer.indent.Close()
	if statement.Else == nil {
		return
	}
	if printer.sawReturn && printer.format.ImproveBranches {
		printer.printNode(statement.Else)
		return
	}
	printer.printElse(statement)
}

func (printer *PrettyPrinter) printElse(statement *ast.ConditionalStatement) {
	printer.appendIndent()
	printer.append(token.ElseKeyword.String())
	printer.appendRune(' ')
	if _, ok := statement.Else.(*ast.ConditionalStatement); !ok {
		printer.appendLineBreak()
		printer.indent.Open()
		defer printer.indent.Close()
	}
	printer.printNode(statement.Else)
}

func (printer *PrettyPrinter) printIfHeader(statement *ast.ConditionalStatement) {
	printer.append(token.IfKeyword.String())
	printer.appendRune(' ')
	printer.printNode(statement)
	printer.appendLineBreak()
}

func (printer *PrettyPrinter) printReturnStatement(statement *ast.ReturnStatement) {
	printer.sawReturn = true
	printer.append(token.ReturnKeyword.String())
	printer.appendRune(' ')
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printYieldStatement(statement *ast.YieldStatement) {
	printer.append(token.YieldKeyword.String())
	printer.appendRune(' ')
	printer.printNode(statement.Value)
}