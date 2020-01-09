package pass

import (
	"fmt"
	"gitlab.com/strict-lang/sdk/pkg/compiler/grammar/tree"
)

type Execution struct {
	target *Pass
	context *Context
}

func NewExecution(target *Pass, unit *tree.TranslationUnit) *Execution {
	return &Execution{
		target: target,
		context: &Context{
			Unit:unit,
		},
	}
}

func (execution *Execution) Run() error {
	orderedPasses, err := execution.orderPendingPasses()
	if err != nil {
		return err
	}
	for _, pass := range orderedPasses {
		pass.Run(execution.context)
	}
	return nil
}

func (execution *Execution) orderPendingPasses() ([]*Pass, error) {
	graph := execution.createDependencyGraph()
	return graph.sortTopologically()
}

func (execution *Execution) createDependencyGraph() *dependencyGraph {
	graph := newDependencyGraph()
	execution.traversePassDependencies(func(pass *Pass, dependency *Pass) {
		graph.insert(pass, dependency)
	})
	return graph
}

type dependencyRelationVisitor func(pass *Pass, dependency *Pass)

func (execution *Execution) traversePassDependencies(
	visitor dependencyRelationVisitor) {

	execution.traversePassDependenciesRecursive(
		execution.target, visitor)
}

func (execution *Execution) traversePassDependenciesRecursive(
	pass *Pass, visitor dependencyRelationVisitor) {

	for _, dependency := range pass.Dependencies {
		visitor(pass, dependency)
		execution.traversePassDependenciesRecursive(pass, visitor)
	}
}

type dependencyGraph struct {
	edges map[*Pass] []*Pass
	elementCount int
}

func newDependencyGraph() *dependencyGraph {
	return &dependencyGraph{}
}

func (graph *dependencyGraph) insert(dependant *Pass, dependency *Pass) {
	current := graph.listDependenciesFor(dependant)
	if dependencies, updated := appendToSet(current, dependency); updated {
		graph.elementCount++
		graph.edges[dependant] = dependencies
	}
}

func appendToSet(set []*Pass, element *Pass) (result []*Pass, updated bool) {
	for _, entry := range set {
		if entry == element {
			return set, false
		}
	}
	return append(set, element), true
}

func (graph *dependencyGraph) listDependenciesFor(pass *Pass) []*Pass {
	if dependencies, isRegistered := graph.edges[pass]; isRegistered {
		return dependencies
	}
	return []*Pass{}
}

func (graph *dependencyGraph) isCircular(node *Pass, dependency *Pass) bool {
	dependencyDependencies := graph.listDependenciesFor(dependency)
	for _, entry := range dependencyDependencies {
		if entry == node {
			return true
		}
	}
	return false
}

func newCircularDependencyError(node *Pass, dependency *Pass) error {
	return fmt.Errorf("circular dependencies: %v, %v", node, dependency)
}

func (graph *dependencyGraph) sortTopologically() ([]*Pass, error) {
	for node, dependencies := range graph.edges {
		for _, dependency := range dependencies {
			if graph.isCircular(node, dependency) {
				return nil, newCircularDependencyError(node, dependency)
			}
		}
	}
	return []*Pass{}, nil
}