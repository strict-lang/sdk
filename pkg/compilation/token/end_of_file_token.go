package token

const (
	EndOfFileTokenName = "eof"
)

var (
	EndOfFile Token = &EndOfFileToken{}
)

type EndOfFileToken struct{}

func (EndOfFileToken) Position() Position {
	return Position{BeginOffset: 0, EndOffset: 0}
}

func (EndOfFileToken) Value() string {
	return ""
}

func (EndOfFileToken) Name() string {
	return EndOfFileTokenName
}

func (EndOfFileToken) Indent() Indent {
	return 0
}

func (EndOfFileToken) String() string {
	return EndOfFileTokenName
}

func IsEndOfFileToken(token Token) bool {
	_, ok := token.(*EndOfFileToken)
	return ok
}
