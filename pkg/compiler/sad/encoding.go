package sad

import (
	"bufio"
	"strconv"
)

type symbolTable struct {
	count int
	table map[string] int
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
	output *bufio.Writer
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
	_, _ = encoding.output.WriteRune(methodBeginKey)
}

func (encoding *encoding) beginClassParameterList() {
	_, _ = encoding.output.WriteRune(classParameterListBegin)
}

func (encoding *encoding) endClassParameterList() {
	_, _ = encoding.output.WriteRune(classParameterListEnd)
}

func (encoding *encoding) beginField() {
	_, _ = encoding.output.WriteRune(fieldBeginKey)
}

func (encoding *encoding) beginClass() {
	_, _ = encoding.output.WriteRune(classBeginKey)
}

const itemSeparator = ';'
const classSeparator = '\n'

func (encoding *encoding) completeItem() {
	_, _ = encoding.output.WriteRune(itemSeparator)
}

func (encoding *encoding) completeClass() {
	_, _ = encoding.output.WriteRune(classSeparator)
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
		class.encodeParameters(encoding)
		encoding.completeItem()
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
	encoding.completeItem()
	for _, parameter := range method.Parameters {
		parameter.encode(encoding)
		encoding.completeItem()
	}
	method.ReturnType.encode(encoding)
	encoding.completeItem()
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
	encoding.completeItem()
}