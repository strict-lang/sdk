package decoder

import "gitlab.com/strict-lang/sdk/sir"

type Decoder interface {
	Decode([]byte) *sir.Unit
}
