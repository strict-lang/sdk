package dependency

import (
	"fmt"
	"io"
)

type DependChain struct {
	dependencies []*Dependency
}

func (chain *DependChain) Length() int {
	return len(chain.dependencies)
}

func (chain *DependChain) Root() *Dependency {
	return chain.dependencies[dependChainRootIndex]
}

func (chain *DependChain) Target() *Dependency {
	lastIndex := len(chain.dependencies) - 1
	return chain.dependencies[lastIndex]
}

func (chain *DependChain) Write(writer io.Writer) error {
	return writeDependencyChainElements(chain.dependencies, writer)
}

func (chain *DependChain) WriteLimited(writer io.Writer, limit int) error {
	if chain.Length() <= limit {
		return writeDependencyChainElements(chain.dependencies, writer)
	}
	limited := chain.listRelevantElements(limit)
	return writeDependencyChainElements(limited, writer)
}

func (chain *DependChain) listRelevantElements(limit int) []*Dependency {
	lastHeadElement := limit / 2
	lastIndex := chain.Length() - 1
	lastTailElement := lastIndex - limit/2
	headElements := chain.dependencies[0:lastHeadElement]
	tailElements := chain.dependencies[lastTailElement:lastIndex]
	return append(headElements, tailElements...)
}

func writeDependencyChainElements(elements []*Dependency, writer io.Writer) error {
	lastIndex := len(elements) - 1
	for _, element := range elements {
		if err := writeDependencyChainElement(element, writer); err != nil {
			return err
		}
	}
	lastElement := elements[lastIndex]
	return writeLastDependencyChainElement(lastElement, writer)
}

func writeDependencyChainElement(element *Dependency, writer io.Writer) error {
	formatted := fmt.Sprintf("%s\nreferences ", element.Id())
	_, err := writer.Write([]byte(formatted))
	return err
}

func writeLastDependencyChainElement(element *Dependency, writer io.Writer) error {
	formatted := fmt.Sprintf("%s", element.Id())
	_, err := writer.Write([]byte(formatted))
	return err
}
