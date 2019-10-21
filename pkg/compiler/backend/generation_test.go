package backend

import (
	"fmt"
)

func (generation *Generation) PrintOutput() {
	fmt.Println(generation.output.String())
}
