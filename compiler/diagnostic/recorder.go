package diagnostic

import "gitlab.com/strict-lang/sdk/compiler/source"

type RecordedEntry struct {
	Kind     *Kind
	Stage    *Stage
	Source   string
	Message  string
	UnitName string
	Offset   source.Offset
}

type Recorder struct {
	entries *[]RecordedEntry
}

func NewRecorder() *Recorder {
	return &Recorder{
		entries: &[]RecordedEntry{},
	}
}

func (recorder *Recorder) Record(entry RecordedEntry) {
	*recorder.entries = append(*recorder.entries, entry)
}

type OffsetToPositionConverter func(source.Offset) source.Position

func (recorder *Recorder) CreateDiagnostics(converter OffsetToPositionConverter) *Diagnostics {
	mappedEntries := make(map[string][]Entry)
	for _, recorded := range *recorder.entries {
		position := converter(recorded.Offset)
		entry := Entry{
			Position: Position{
				LineIndex: position.Line.Index,
				Column:    position.Column,
			},
			UnitName: recorded.UnitName,
			Kind:     recorded.Kind,
			Source:   recorded.Source,
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
