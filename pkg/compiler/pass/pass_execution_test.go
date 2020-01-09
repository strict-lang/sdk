package pass

import (
	"gitlab.com/strict-lang/sdk/pkg/compiler/isolate"
	"testing"
)

type testPass struct {
	output *sharedSlice
	value int
	dependencies []Pass
}

func (pass *testPass) Run(context *Context) {
	pass.output.value = append(pass.output.value, pass.value)
}

func (pass *testPass) Dependencies(*isolate.Isolate) Set {
	return pass.dependencies
}

type sharedSlice struct {
	value []int
}

func TestExecution_Run(testing *testing.T) {
	var result sharedSlice
	first := &testPass{
		output: &result,
		value: 1,
	}
	second := &testPass{
		output: &result,
		value: 2,
		dependencies: Set{first},
	}
	third := &testPass{
		output: &result,
		value: 3,
		dependencies: Set{first, second},
	}

	execution := NewExecution(third, &Context{
		Isolate: isolate.SingleThreaded(),
	})
	if err := execution.Run(); err != nil {
		testing.Error(err)
	}
	if !isSliceEqual(result.value, []int{1, 2, 3}) {
		testing.Errorf("wrong order of execution: %v", result)
	}
}

func isSliceEqual(left []int, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for index, value := range left {
		if right[index] != value {
			return false
		}
	}
	return true
}
