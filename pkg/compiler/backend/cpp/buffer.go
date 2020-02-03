package cpp

type OutputStream interface {
	WriteLine(string)
	Write(string)
	WriteFormatted(string, ...interface{})
}

type FormattedOutputStream interface {
	OutputStream
	Indent()
	IncreaseIndent()
	DecreaseIndent()
}
