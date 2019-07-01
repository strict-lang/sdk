package diagnostic

type Recorder struct {
	entries map[string][]Entry
}

func NewRecorder() *Recorder {
	return nil
}

func (recorder *Recorder) Record(entry Entry) {
	recorder.appendToKind(entry.Kind, entry)
}

func (recorder *Recorder) entriesOfKind(kind *Kind) ([]Entry, bool) {
	entries, ok := recorder.entries[kind.Name]
	return entries, ok
}

func (recorder *Recorder) appendToKind(kind *Kind, entry Entry) {
	recorder.entries[kind.Name] = append(recorder.entries[kind.Name], entry)
}

func (recorder *Recorder) Entries(kind *Kind) ([]Entry, bool) {
	return recorder.entriesOfKind(kind)
}

func (recorder *Recorder) AllEntries() []Entry {
	var all []Entry
	for _, entries := range recorder.entries {
		all = append(all, entries...)
	}
	return all
}

func (recorder *Recorder) PrintAllEntries(printer Printer) {
	for _, entries := range recorder.entries {
		for _, entry := range entries {
			entry.PrintColored(printer)
		}
	}
}