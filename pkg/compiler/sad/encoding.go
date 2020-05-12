package sad

import (
	"strconv"
	"strings"
)

func Encode(tree *Tree) string {
	encoding := newEncoding(tree)
	return encoding.generate()
}

type symbolTable struct {
	count int
	table map[string]int
	ordered []string
}

func newSymbolTable() symbolTable {
	return symbolTable{
		count: 0,
		table: map[string]int{},
	}
}

func (symbols *symbolTable) register(name string) int {
	if index, ok := symbols.table[name]; ok {
		return index
	}
	nextIndex := symbols.count
	symbols.table[name] = nextIndex
	symbols.ordered = append(symbols.ordered, name)
	symbols.count++
	return nextIndex
}

type encoding struct {
	symbols symbolTable
	output  *strings.Builder
	tree    *Tree
}

func newEncoding(tree *Tree) *encoding {
	return &encoding{
		symbols: newSymbolTable(),
		output:  &strings.Builder{},
		tree:    tree,
	}
}

func (encoding *encoding) generate() string {
	for _, class := range encoding.tree.Classes {
		class.encode(encoding)
	}
	tail := encoding.output.String()
	encoding.output = &strings.Builder{}
	encoding.writeSymbols()
	encoding.completeSymbolList()
	return encoding.output.String() + tail
}

func (encoding *encoding) writeSymbols() {
	for index, symbol := range encoding.symbols.ordered {
		encoding.output.WriteString(symbol)
		if index != len(encoding.symbols.ordered) - 1 {
			encoding.completeSymbol()
		}
	}
}

const symbolListSeparator = '\n'
const symbolSeparator = ';'

func (encoding *encoding) completeSymbolList() {
	encoding.writeRune(symbolListSeparator)
}

func (encoding *encoding) completeSymbol() {
	encoding.writeRune(symbolSeparator)
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
	for index, symbol := range encoding.symbols.ordered {
		_, _ = encoding.output.WriteString(symbol)
		if index != len(encoding.symbols.ordered) - 1 {
			_, _ = encoding.output.WriteRune(itemSeparator)
		}
	}
}

func (class *Class) encode(encoding *encoding) {
	encoding.beginClass()
	encoding.writeSymbol(class.Name)
	encoding.completeItem()
	class.maybeEncodeParameters(encoding)
	class.encodeTraits(encoding)
	encoding.completeClassItem()
	class.encodeItems(encoding)
	encoding.completeClass()
}

func (class *Class) encodeTraits(encoding *encoding) {
	for index, trait := range class.Traits {
		trait.encode(encoding)
		if index != len(class.Traits) - 1 {
			encoding.completeItem()
		}
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
	for index, argument := range name.Arguments {
		argument.encode(encoding)
		if index != len(name.Arguments) - 1 {
			encoding.completeItem()
		}
	}
	encoding.endClassParameterList()
}

func (method *Method) encode(encoding *encoding) {
	encoding.beginMethod()
	encoding.writeSymbol(method.Name)
	method.encodeParameters(encoding)
	method.ReturnType.encode(encoding)
	encoding.completeClassItem()
}

func (method *Method) encodeParameters(encoding *encoding) {
	encoding.beginParameterList()
	for index, parameter := range method.Parameters {
		parameter.encode(encoding)
		if index != len(method.Parameters) - 1 {
			encoding.completeParameter()
		}
	}
	encoding.endParameterList()
}

func (parameter *Parameter) encode(encoding *encoding) {
	if len(parameter.Label) != 0 {
		encoding.writeSymbol(parameter.Label)
		encoding.completeItem()
	}
	encoding.writeSymbol(parameter.Name)
	encoding.completeItem()
	parameter.Class.encode(encoding)
}

func (field *Field) encode(encoding *encoding) {
	encoding.beginField()
	encoding.writeSymbol(field.Name)
	encoding.completeItem()
	field.Class.encode(encoding)
	encoding.completeClassItem()
}
