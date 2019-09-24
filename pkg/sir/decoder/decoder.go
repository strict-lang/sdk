package decoder

import (
	sir2 "gitlab.com/strict-lang/sdk/pkg/sir"
)

type Decoder interface {
	Decode([]byte) *sir2.Unit
}
