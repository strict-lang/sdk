package sad

import (
	"bufio"
	"strconv"
)

type symbolTable struct {
	count int
	table map[string]int
}

func (symbols *symbolTable) register(name string) int {
	if index, ok := symbols.table[name]; ok {
		return index
	}
	nextIndex := symbols.count
	symbols.table[name] = nextIndex
	symbols.count++
	return nextIndex
}

type encoding struct {
	symbols symbolTable
	output  *bufio.Writer
}

func (encoding *encoding) writeSymbol(name string) {
	index := encoding.symbols.register(name)
	_, _ = encoding.output.WriteString(strconv.Itoa(index))
}

func (encoding *encoding) writeRune(value rune) {
	_, _ = encoding.output.WriteRune(value)
}

const methodBeginKey = 'm'
const fieldBeginKey = 'f'
const classBeginKey = 'c'
const symbolTableBeginKey = 's'

const classParameterListBegin = '<'
const classParameterListEnd = '>'

func (encoding *encoding) beginMethod() {
	encoding.writeRune(methodBeginKey)
}

func (encoding *encoding) beginClassParameterList() {
	encoding.writeRune(classParameterListBegin)
}

func (encoding *encoding) endClassParameterList() {
	encoding.writeRune(classParameterListEnd)
}

func (encoding *encoding) beginField() {
	encoding.writeRune(fieldBeginKey)
}

func (encoding *encoding) beginClass() {
	encoding.writeRune(classBeginKey)
}

const parameterListBegin = '('
const parameterListEnd = ')'
const parameterSeparator = ','

func (encoding *encoding) completeParameter() {
	encoding.writeRune(parameterSeparator)
}

func (encoding *encoding) beginParameterList() {
	encoding.writeRune(parameterListBegin)
}

func (encoding *encoding) endParameterList() {
	encoding.writeRune(parameterListEnd)
}

const itemSeparator = '.'
const classItemSeparator = ';'
const classSeparator = '\n'

func (encoding *encoding) completeItem() {
	encoding.writeRune(itemSeparator)
}

func (encoding *encoding) completeClassItem() {
	encoding.writeRune(classItemSeparator)
}

func (encoding *encoding) completeClass() {
	encoding.writeRune(classSeparator)
}

func (encoding *encoding) writeSymbolTable() {
	_, _ = encoding.output.WriteRune(symbolTableBeginKey)
	for symbol := range encoding.symbols.table {
		_, _ = encoding.output.WriteString(symbol)
		_, _ = encoding.output.WriteRune(itemSeparator)
	}
}

func (class *Class) encode(encoding *encoding) {
	encoding.beginClass()
	encoding.writeSymbol(class.Name)
	class.maybeEncodeParameters(encoding)
	class.encodeTraits(encoding)
	class.encodeItems(encoding)
	encoding.completeClass()
}

func (class *Class) encodeTraits(encoding *encoding) {
	for _, trait := range class.Traits {
		trait.encode(encoding)
		encoding.completeItem()
	}
}

func (class *Class) encodeItems(encoding *encoding) {
	for _, method := range class.Methods {
		method.encode(encoding)
	}
	for _, field := range class.Fields {
		field.encode(encoding)
	}
}

func (class *Class) maybeEncodeParameters(encoding *encoding) {
	if len(class.Parameters) != 0 {
		encoding.beginParameterList()
		class.encodeParameters(encoding)
		encoding.endParameterList()
	}
}

func (class *Class) encodeParameters(encoding *encoding) {
	encoding.beginClassParameterList()
	for _, parameter := range class.Parameters {
		if parameter.Wildcard {
			encoding.writeSymbol("*")
		}
		encoding.writeSymbol(parameter.Class.Name)
	}
	encoding.endClassParameterList()
}

func (name *ClassName) encode(encoding *encoding) {
	if name.Wildcard {
		encoding.writeSymbol("*")
		if len(name.Name) == 0 {
			return
		}
	}
	encoding.writeSymbol(name.Name)
	if len(name.Arguments) != 0 {
		name.encodeArguments(encoding)
	}
}

func (name *ClassName) encodeArguments(encoding *encoding) {
	encoding.beginClassParameterList()
	for _, argument := range name.Arguments {
		argument.encode(encoding)
		encoding.completeItem()
	}
	encoding.endClassParameterList()
}

func (method *Method) encode(encoding *encoding) {
	encoding.beginMethod()
	encoding.writeSymbol(method.Name)
	method.encodeParameters(encoding)
	encoding.completeItem()
	method.ReturnType.encode(encoding)
	encoding.completeClassItem()
}

func (method *Method) encodeParameters(encoding *encoding) {
	encoding.beginParameterList()
	for _, parameter := range method.Parameters {
		parameter.encode(encoding)
		encoding.completeParameter()
	}
	encoding.endParameterList()
}

func (parameter *Parameter) encode(encoding *encoding) {
	encoding.writeRune('(')
	if len(parameter.Label) != 0 {
		encoding.writeSymbol(parameter.Label)
		encoding.completeItem()
	}
	encoding.writeSymbol(parameter.Name)
	encoding.completeItem()
	parameter.Class.encode(encoding)
	encoding.completeItem()
	encoding.writeRune(')')
}

func (field *Field) encode(encoding *encoding) {
	encoding.beginField()
	encoding.writeSymbol(field.Name)
	encoding.completeItem()
	field.Class.encode(encoding)
	encoding.completeClassItem()
}
