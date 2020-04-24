package diagnostic

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
)

type RecordedEntry struct {
	Kind     *Kind
	Stage    *Stage
	Message  string
	UnitName string
	Error    *RichError
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
	if entry.Message == "" {
		entry.Message = entry.Error.Error.Name()
	}
	*recorder.entries = append(*recorder.entries, entry)
}

type OffsetConversionFunction func(input.Offset) input.Position

func (recorder *Bag) CreateDiagnostics(converter OffsetConversionFunction) *Diagnostics {
	mappedEntries := make(map[string][]Entry)
	for _, recorded := range *recorder.entries {
		position := converter(recorded.Position.Begin())
		entry := Entry{
			Position: Position{
				Begin: position,
				End:   converter(recorded.Position.End()),
			},
			Source:   position.Line.Text,
			UnitName: recorded.UnitName,
			Kind:     recorded.Kind,
			Message:  recorded.Message,
			Stage:    recorded.Stage,
			Error:    recorded.Error,
		}
		kindName := entry.Kind.Name
		current := mappedEntries[kindName]
		mappedEntries[kindName] = append(current, entry)
	}
	return &Diagnostics{
		entries: mappedEntries,
	}
}
