package linemap

import "github.com/BenjaminNitschke/Strict/compiler/source"

type Builder struct {
}

func (builder *Builder) Append(length source.Offset) {

}

func (builder *Builder) NewLinemap() *Linemap {
	return nil
}

func NewBuilder() *Builder {
	return &Builder{}
}
