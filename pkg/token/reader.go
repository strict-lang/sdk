package token

type Reader interface {
	Pull() Token
	Peek() Token
	Last() Token
}
