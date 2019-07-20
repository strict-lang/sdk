package format

const (
	WindowsEndOfLine = "\r\n"
	UnixEndOfLine    = "\n"
)

type Format struct {
	EndOfLine       string
	IndentWriter    IndentWriter
	TabWidth        int
	LineLengthLimit int
	ImproveBranches bool
}
