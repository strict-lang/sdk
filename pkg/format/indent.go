package format

type Indent struct {
	Level      int
	Continuous bool
}

func (indent *Indent) OpenContinuous() {
	indent.Continuous = true
	indent.Level++
}

func (indent *Indent) CloseContinuous() {
	indent.Continuous = false
	indent.Level--
}

func (indent *Indent) Open() {
	indent.Level++
}

func (indent *Indent) Close() {
	indent.Level--
}

type IndentWriter interface {
	Write(indent Indent, writer Writer)
}

type TabIndentWriter struct{}

func (TabIndentWriter) Write(indent Indent, writer Writer) {
	writeRepeatedRune('\t', indent.Level, writer)
}

type SimpleSpaceIndentWriter struct {
	SpacesPerLevel int
}

func (spaceWriter SimpleSpaceIndentWriter) Write(indent Indent, writer Writer) {
	spaces := spaceWriter.SpacesPerLevel * indent.Level
	writeRepeatedRune(' ', spaces, writer)
}

type ComplexSpaceIndentWriter struct {
	SpacesPerLevel   int
	ContinuousIndent int
}

func (complexWriter ComplexSpaceIndentWriter) Write(indent Indent, writer Writer) {
	spaces := complexWriter.spacesForIndent(indent)
	writeRepeatedRune(' ', spaces, writer)
}

func (complexWriter ComplexSpaceIndentWriter) spacesForIndent(indent Indent) int {
	if indent.Continuous {
		return (indent.Level-1)*complexWriter.SpacesPerLevel + complexWriter.ContinuousIndent
	}
	return indent.Level * complexWriter.SpacesPerLevel
}

func writeRepeatedRune(value rune, count int, writer Writer) {
	for index := 0; index < count; index++ {
		WriteRune(value)
	}
}
