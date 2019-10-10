package diagnostic

import (
	 "gitlab.com/strict-lang/sdk/pkg/compiler/input"
)

type RecordedEntry struct {
	Kind     *Kind
	Stage    *Stage
	Message  string
	UnitName string
	Position RecordedPosition
}

type Bag struct {
	entries *[]RecordedEntry
}

type RecordedPosition interface {
	Begin() input.Offset
	End() input.Offset
}

func NewBag() *Bag {
	return &Bag{
		entries: &[]RecordedEntry{},
	}
}

func (recorder *Bag) Record(entry RecordedEntry) {
	*recorder.entries = append(*recorder.entries, entry)
}

type OffsetConversionFunction func(input.Offset) input.Position

func (recorder *Bag) CreateDiagnostics(converter OffsetConversionFunction) *Diagnostics {
	mappedEntries := make(map[string][]Entry)
	for _, recorded := range *recorder.entries {
		position := converter(recorded.Position.Begin())
		entry := Entry{
			Position: Position{
				LineIndex: position.Line.Index,
				Column:    position.Column,
			},
			UnitName: recorded.UnitName,
			Kind:     recorded.Kind,
			Message:  recorded.Message,
			Stage:    recorded.Stage,
		}
		kindName := entry.Kind.Name
		current := mappedEntries[kindName]
		mappedEntries[kindName] = append(current, entry)
	}
	return &Diagnostics{
		entries: mappedEntries,
	}
}
