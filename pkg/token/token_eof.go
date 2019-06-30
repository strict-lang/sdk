package token

const (
	EndOfFileTokenName = "EndOfFile"
)

var (
	EndOfFile Token = &EndOfFileToken{}
)
type EndOfFileToken struct {}

func (EndOfFileToken) Position() Position {
	return Position{Begin: 0, End: 0,}
}

func (EndOfFileToken) Value() string {
	return ""
}

func (EndOfFileToken) Name() string {
	return EndOfFileTokenName
}

func (EndOfFileToken) IsOperator() bool {
	return false
}

func (EndOfFileToken) IsKeyword() bool {
	return false
}

func (EndOfFileToken) IsValid() bool {
	return true
}

func (EndOfFileToken) IsLiteral() bool {
	return false
}
