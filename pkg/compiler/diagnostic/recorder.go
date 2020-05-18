package diagnostic

import (
	"github.com/strict-lang/sdk/pkg/compiler/input"
	"github.com/strict-lang/sdk/pkg/compiler/input/linemap"
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

func ConvertWithLineMap(lineMap *linemap.LineMap) OffsetConversionFunction {
	return lineMap.PositionAtOffset
}

func (recorder *Bag) CreateDiagnostics(converter OffsetConversionFunction) *Diagnostics {
	var entries []Entry
	for _, recorded := range *recorder.entries {
		entry := translateEntry(converter, recorded)
		entries = append(entries, entry)
	}
	return &Diagnostics{entries: entries}
}

func translateEntry(
	converter OffsetConversionFunction,
	recorded RecordedEntry) Entry {

	position := converter(recorded.Position.Begin())
	return Entry{
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
}