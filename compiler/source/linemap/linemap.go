package linemap

import "github.com/BenjaminNitschke/Strict/compiler/source"

type Linemap struct {
	Offsets []source.Position
}

func (lines *Linemap) LineAtOffset(offset source.Offset) source.LineIndex {
	return 0
}

func (lines *Linemap) OffsetAtLine(index source.LineIndex) source.Offset {
	return 0
}

func (lines *Linemap) PositionAtOffset(offset source.Offset) source.Position {
	return source.Position{}
}

func (lines *Linemap) PositionAtLine(index source.LineIndex) source.Position {
	return source.Position{}
}
