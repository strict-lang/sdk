package source

type Reader interface {
	Pull() Char
	Peek() Char
	Last() Char
	Index() Offset
	Skip(count int)
	IsExhausted() bool
}
