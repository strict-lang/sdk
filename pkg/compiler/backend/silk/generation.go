package silk

import "github.com/strict-lang/sdk/pkg/silk"

type Generation struct {
}

func (generation *Generation) EmitInstruction(instruction *silk.Instruction) {

}

func (generation *Generation) selectLoadTarget(item *Item) (StorageLocation, error) {
	return nil, nil
}
