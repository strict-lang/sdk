package diagnostic

type Recorder interface {
	Record(Entry)
}
