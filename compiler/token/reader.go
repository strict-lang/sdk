package token

// Reader interface represents a peekable stream of tokens. It is implemented
// by the scanner and provides a rather narrow interface. Other modules like
// the parser, only depend on the token.Reader and not the scanner, this allows
// multiple reader implementations to be made, which makes testing easier.
type Reader interface {
	// Pull pulls the next token from the stream. If there is no next token, an
	// EndOfFile token is returned. Subsequent calls to Pull will never return
	// the same token. The most recently pulled token is returned by Last().
	Pull() Token
	// Peek peeks the next token in the stream without modifying it. If there is
	// no next token, an EndOfFile token is returned. Subsequent calls to Peek will
	// always return the same token until Pull() is called.
	Peek() Token
	// Last returns the most recently pulled token.
	Last() Token
}
