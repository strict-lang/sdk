package source

type Reader interface {
	Pull() Char
	Peek() Char
	Last() Char
	Skip(count int)
	IsExhausted() bool
}
