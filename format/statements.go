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
	printer.printNode(statement.Consequence)
	printer.indent.Close()
	if !statement.HasAlternative() {
		return
	}
	if printer.sawReturn && printer.format.ImproveBranches {
		printer.printNode(statement.Alternative)
		return
	}
	printer.printElse(statement)
}

func (printer *PrettyPrinter) printElse(statement *ast.ConditionalStatement) {
	printer.appendIndent()
	printer.append("else ")
	if _, ok := statement.Alternative.(*ast.ConditionalStatement); !ok {
		printer.appendLineBreak()
		printer.indent.Open()
		defer printer.indent.Close()
	}
	printer.printNode(statement.Alternative)
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

func (printer *PrettyPrinter) printForEachLoopStatement(loop *ast.ForEachLoopStatement) {
	printer.appendFormatted(
		"for %s in ", loop.Field.Value)

	printer.printNode(loop.Enumeration)
	printer.append(" do")
	printer.appendLineBreak()
	printer.indent.Open()
	defer printer.indent.Close()
	printer.printNode(loop.Body)
}

func (printer *PrettyPrinter) printRangedLoopStatement(loop *ast.RangedLoopStatement) {
	printer.appendFormatted(
		"for %s from ", loop.ValueField.Value)

	printer.printNode(loop.InitialValue)
	printer.appendFormatted(" to ")
	printer.printNode(loop.EndValue)
	printer.append(" do")
	printer.appendLineBreak()
	printer.indent.Open()
	defer printer.indent.Close()
	printer.printNode(loop.Body)
}
