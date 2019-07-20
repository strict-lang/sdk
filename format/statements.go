package format

import (
	"gitlab.com/strict-lang/sdk/compiler/ast"
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
	printer.append("else ")
	if _, ok := statement.Else.(*ast.ConditionalStatement); !ok {
		printer.appendLineBreak()
		printer.indent.Open()
		defer printer.indent.Close()
	}
	printer.printNode(statement.Else)
}

func (printer *PrettyPrinter) printIfHeader(statement *ast.ConditionalStatement) {
	printer.append("if ")
	printer.printNode(statement.Condition)
	printer.append(" do")
	printer.appendLineBreak()
}

func (printer *PrettyPrinter) printReturnStatement(statement *ast.ReturnStatement) {
	printer.sawReturn = true
	printer.append("return ")
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printYieldStatement(statement *ast.YieldStatement) {
	printer.append("yield ")
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printAssignStatement(statement *ast.AssignStatement) {
	printer.printNode(statement.Target)
	printer.appendFormatted(" %s ", statement.Operator.String())
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printForeachLoopStatement(loop *ast.ForeachLoopStatement) {
	printer.appendFormatted(
		"for %s in ", loop.Field.Value)

	printer.printNode(loop.Target)
	printer.append(" do")
	printer.appendLineBreak()
	printer.indent.Open()
	defer printer.indent.Close()
	printer.printNode(loop.Body)
}

func (printer *PrettyPrinter) printFromToLoopStatement(loop *ast.FromToLoopStatement) {
	printer.appendFormatted(
		"for %s from ", loop.Field.Value)

	printer.printNode(loop.From)
	printer.appendFormatted(" to ")
	printer.printNode(loop.To)
	printer.append(" do")
	printer.appendLineBreak()
	printer.indent.Open()
	defer printer.indent.Close()
	printer.printNode(loop.Body)
}
