package silk

type Field struct {
	Name string
	Index int
	Type Type
}

func (field *Field) Matches(target StorageLocation) bool {
	return false
}
