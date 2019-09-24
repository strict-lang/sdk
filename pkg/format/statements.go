package format

import (
	syntaxtree2 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
)

func (printer *PrettyPrinter) printBlockStatement(block *syntaxtree2.BlockStatement) {
	for _, statement := range block.Children {
		printer.appendIndent()
		printer.printNode(statement)
		printer.appendLineBreak()
	}
}

func (printer *PrettyPrinter) printExpressionStatement(statement *syntaxtree2.ExpressionStatement) {
	printer.printNode(statement.Expression)
}

func (printer *PrettyPrinter) printConditionalStatement(statement *syntaxtree2.ConditionalStatement) {
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

func (printer *PrettyPrinter) printElse(statement *syntaxtree2.ConditionalStatement) {
	printer.appendIndent()
	printer.append("else ")
	if _, ok := statement.Alternative.(*syntaxtree2.ConditionalStatement); !ok {
		printer.appendLineBreak()
		printer.indent.Open()
		defer printer.indent.Close()
	}
	printer.printNode(statement.Alternative)
}

func (printer *PrettyPrinter) printIfHeader(statement *syntaxtree2.ConditionalStatement) {
	printer.append("if ")
	printer.printNode(statement.Condition)
	printer.append(" do")
	printer.appendLineBreak()
}

func (printer *PrettyPrinter) printReturnStatement(statement *syntaxtree2.ReturnStatement) {
	printer.sawReturn = true
	printer.append("return ")
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printYieldStatement(statement *syntaxtree2.YieldStatement) {
	printer.append("yield ")
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printAssignStatement(statement *syntaxtree2.AssignStatement) {
	printer.printNode(statement.Target)
	printer.appendFormatted(" %s ", statement.Operator.String())
	printer.printNode(statement.Value)
}

func (printer *PrettyPrinter) printForEachLoopStatement(loop *syntaxtree2.ForEachLoopStatement) {
	printer.appendFormatted(
		"for %s in ", loop.Field.Value)

	printer.printNode(loop.Sequence)
	printer.append(" do")
	printer.appendLineBreak()
	printer.indent.Open()
	defer printer.indent.Close()
	printer.printNode(loop.Body)
}

func (printer *PrettyPrinter) printRangedLoopStatement(loop *syntaxtree2.RangedLoopStatement) {
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
