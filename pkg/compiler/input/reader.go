package input

type Reader interface {
	Pull() Char
	Peek() Char
	Current() Char
	Index() Offset
	Skip(count int)
	IsExhausted() bool
}
