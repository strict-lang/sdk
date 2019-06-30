package diagnostic

type Recorder interface {
	Record(Entry)
}

func NewRecorder() Recorder {
	return nil
}