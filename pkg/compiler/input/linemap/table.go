package linemap

type Table struct {
	maps map[string] *LineMap
}

func NewTable(maps map[string] *LineMap) *Table {
	return &Table{maps: maps}
}

func NewEmptyTable() *Table {
	return &Table{maps: map[string]*LineMap{}}
}

func (table *Table) Insert(name string, lineMap *LineMap) {
	table.maps[name] = lineMap
}

func (table *Table) Lookup(name string) (*LineMap, bool) {
	lineMap, ok := table.maps[name]
	return lineMap, ok
}
