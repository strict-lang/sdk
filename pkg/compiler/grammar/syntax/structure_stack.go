package syntax

import (
	"errors"
	"fmt"
	"strict.dev/sdk/pkg/compiler/grammar/tree"
	"strict.dev/sdk/pkg/compiler/input"
	"strings"
)

type structureStack struct {
	elements            []structureStackElement
	history             strings.Builder
	historyIndent       int
	shouldRecordHistory bool
}

func newRecordingStructureStack() *structureStack {
	return &structureStack{
		elements:            nil,
		history:             strings.Builder{},
		historyIndent:       0,
		shouldRecordHistory: true,
	}
}

type structureStackElement struct {
	nodeKind    tree.NodeKind
	beginOffset input.Offset
}

func (stack *structureStack) listRemainingElementsOrdered() []structureStackElement {
	length := len(stack.elements)
	reversed := make([]structureStackElement, length)
	for index, element := range stack.elements {
		reversed[length-1-index] = element
	}
	return reversed
}

func (stack *structureStack) isEmpty() bool {
	return len(stack.elements) == 0
}

func (stack *structureStack) createHistoryDump() string {
	return stack.history.String()
}

var errEmptyStructureStack = errors.New("structureStack is empty")

func (stack *structureStack) peek() structureStackElement {
	index := len(stack.elements) - 1
	if index < 0 {
		panic(errEmptyStructureStack)
	}
	return stack.elements[index]
}

func (stack *structureStack) push(element structureStackElement) {
	stack.elements = append(stack.elements, element)
	if stack.shouldRecordHistory {
		stack.writePushedElementToHistory(element)
	}
}

const structureStackHistoryIndent = '.'
const structureStackHistoryLineBegin = " | "

func (stack *structureStack) writePushedElementToHistory(element structureStackElement) {
	indent := repeatRune(structureStackHistoryIndent, stack.historyIndent)
	stack.history.WriteString(structureStackHistoryLineBegin)
	message := fmt.Sprintf("%s >> %s at offset %d\n", indent, element.nodeKind, element.beginOffset)
	_, _ = stack.history.WriteString(message)
	stack.historyIndent++
}

func (stack *structureStack) pop() (structureStackElement, error) {
	index := len(stack.elements) - 1
	if index < 0 {
		return structureStackElement{}, errEmptyStructureStack
	}
	element := stack.elements[index]
	stack.elements = stack.elements[0:index]
	if stack.shouldRecordHistory {
		stack.writePoppedElementToHistory(element)
	}
	return element, nil
}

func (stack *structureStack) writePoppedElementToHistory(element structureStackElement) {
	stack.historyIndent--
	indent := repeatRune(structureStackHistoryIndent, stack.historyIndent)
	stack.history.WriteString(structureStackHistoryLineBegin)
	message := fmt.Sprintf("%s << %s\n", indent, element.nodeKind)
	_, _ = stack.history.WriteString(message)
}

func repeatRune(value rune, count int) string {
	runes := make([]rune, count)
	for index := 0; index < count; index++ {
		runes[index] = value
	}
	return string(runes)
}
