package format

import "gitlab.com/strict-lang/sdk/compiler/ast"

func (printer *PrettyPrinter) printMethod(method *ast.Method) {
	printer.appendFormatted(
		"method %s %s(", method.Type.FullName(), method.Name.Value)

	printer.writeMethodParameters(method)
	printer.appendLineBreak()
	printer.indent.Open()
	defer printer.indent.Close()
	printer.printNode(method.Body)
}

func (printer *PrettyPrinter) writeMethodParameters(method *ast.Method) {
	printer.indent.OpenContinuous()
	parameters, combinedLength := printer.recordAllParameters(method)
	lengthOfSpaces := len(parameters) * 2
	totalLineLength := combinedLength + printer.lineLength + lengthOfSpaces
	if totalLineLength >= printer.format.LineLengthLimit {
		printer.writeLongParameterList(parameters)
	} else {
		printer.writeShortParameterList(parameters)
	}
}

func (printer *PrettyPrinter) writeLongParameterList(parameters []string) {
	printer.appendLineBreak()
	printer.appendIndent()
	for index, parameter := range parameters {
		if index != 0 {
			printer.appendRune(',')
			printer.appendLineBreak()
			printer.appendIndent()
		}
		printer.append(parameter)
	}
	printer.appendLineBreak()
	printer.indent.CloseContinuous()
	printer.appendIndent()
	printer.appendRune(')')
}

func (printer *PrettyPrinter) writeShortParameterList(parameters []string) {
	for index, parameter := range parameters {
		if index != 0 {
			printer.append(", ")
		}
		printer.append(parameter)
	}
	printer.indent.CloseContinuous()
	printer.appendRune(')')
}

func (printer *PrettyPrinter) recordAllParameters(
	call *ast.Method) (parameters []string, combinedLength int) {

	for _, parameter := range call.Parameters {
		recorded := printer.recordParameter(parameter)
		parameters = append(parameters, recorded)
		combinedLength += len(recorded)
	}
	return
}

func (printer *PrettyPrinter) recordParameter(parameter ast.Parameter) string {
	buffer := NewStringWriter()
	oldWriter := printer.swapWriter(buffer)
	defer printer.setWriter(oldWriter)
	if isTypeNamedParameter(parameter) {
		printer.append(parameter.Type.FullName())
	} else {
		printer.appendFormatted("%s %s",
			parameter.Type.FullName(), parameter.Name.Value)
	}
	return buffer.String()
}

func isTypeNamedParameter(parameter ast.Parameter) bool {
	return parameter.Name.Value == parameter.Type.NonGenericName()
}
