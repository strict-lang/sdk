package format

import (
	 "gitlab.com/strict-lang/sdk/pkg/compilation/syntaxtree"
	 "gitlab.com/strict-lang/sdk/pkg/compilation/token"
)

func (printer *PrettyPrinter) printIdentifier(identifier *syntaxtree.Identifier) {
	printer.append(identifier.Value)
}

func (printer *PrettyPrinter) printNumberLiteral(number *syntaxtree.NumberLiteral) {
	printer.append(number.Value)
}

func (printer *PrettyPrinter) printStringLiteral(literal *syntaxtree.StringLiteral) {
	printer.appendRune('"')
	printer.append(literal.Value)
	printer.appendRune('"')
}

func (printer *PrettyPrinter) printUnaryExpression(expression *syntaxtree.UnaryExpression) {
	operatorName := keywordNameOrOperator(expression.Operator)
	printer.append(operatorName)
	printer.printNode(expression.Operand)
}

func (printer *PrettyPrinter) printBinaryExpression(expression *syntaxtree.BinaryExpression) {
	printer.printNode(expression.LeftOperand)
	printer.appendRune(' ')

	operatorName := keywordNameOrOperator(expression.Operator)
	printer.append(operatorName)
	printer.appendRune(' ')

	printer.printNode(expression.RightOperand)
}

func keywordNameOrOperator(operator token.Operator) string {
	if keyword, ok := token.KeywordValueOfOperator(operator); ok {
		return keyword.String()
	}
	return operator.String()
}

func (printer *PrettyPrinter) printMethodCall(call *syntaxtree.CallExpression) {
	printer.printNode(call.Method)
	printer.appendRune('(')
	printer.indent.OpenContinuous()
	printer.writeArguments(call)
}

func (printer *PrettyPrinter) writeArguments(call *syntaxtree.CallExpression) {
	arguments, combinedLength := printer.recordAllArguments(call)
	lengthOfSpaces := len(arguments) * 2 // Most arguments have the ', ' chars.
	totalLineLength := combinedLength + printer.lineLength + lengthOfSpaces
	if totalLineLength >= printer.format.LineLengthLimit {
		printer.writeLongArgumentList(arguments)
	} else {
		printer.writeShortArgumentList(arguments)
	}
}

func (printer *PrettyPrinter) writeLongArgumentList(arguments []string) {
	printer.appendLineBreak()
	printer.appendIndent()
	for index, argument := range arguments {
		if index != 0 {
			printer.appendRune(',')
			printer.appendLineBreak()
			printer.appendIndent()
		}
		printer.append(argument)
	}
	printer.appendLineBreak()
	printer.indent.CloseContinuous()
	printer.appendIndent()
	printer.appendRune(')')
}

func (printer *PrettyPrinter) writeShortArgumentList(arguments []string) {
	for index, argument := range arguments {
		if index != 0 {
			printer.append(", ")
		}
		printer.append(argument)
	}
	printer.indent.CloseContinuous()
	printer.appendRune(')')
}

func (printer *PrettyPrinter) recordAllArguments(
	call *syntaxtree.CallExpression) (arguments []string, combinedLength int) {

	for _, argument := range call.Arguments {
		recorded := printer.recordArgument(argument)
		arguments = append(arguments, recorded)
		combinedLength += len(recorded)
	}
	return arguments, combinedLength
}

func (printer *PrettyPrinter) recordArgument(node syntaxtree.Node) string {
	buffer := NewStringWriter()
	oldWriter := printer.swapWriter(buffer)
	defer printer.setWriter(oldWriter)

	printer.printNode(node)
	return buffer.String()
}

func (printer *PrettyPrinter) printSelectorExpression(selector *syntaxtree.SelectExpression) {
	printer.printNode(selector.Target)
	printer.appendRune('.')
	printer.printNode(selector.Selection)
}
