package symbol

type Reference int

type Table struct {
	symbols []Symbol
}

func (table *Table) Lookup(reference Reference) (Symbol, bool) {
	if table.isValidReference(reference) {
		return table.symbols[int(reference)], true
	}
	return nil, false
}

func (table *Table) isValidReference(reference Reference) bool {
	return reference > 0  && int(reference) < len(table.symbols)
}

type TableBuilder struct {
	 symbols []Symbol
}

func (builder* TableBuilder) Insert(symbol Symbol) Reference {
	builder.symbols = append(builder.symbols, symbol)
	assignedIndex := len(builder.symbols) - 1
	return Reference(assignedIndex)
}

func (builder* TableBuilder) CreateTable() *Table {
	copiedSymbols := make([]Symbol, len(builder.symbols))
	copy(copiedSymbols, builder.symbols)
	return &Table{symbols: copiedSymbols}
}