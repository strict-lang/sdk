package pass

import (
	"errors"
	"fmt"
)

type Execution struct {
	target Pass
	context *Context
}

func NewExecution(target Pass, context *Context) *Execution {
	return &Execution{
		target: target,
		context: context,
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

func (execution *Execution) orderPendingPasses() ([]Pass, error) {
  order, err := execution.createDependencyOrder()
  if err != nil {
  	return nil, err
  }
	return order.compute()
}

func (execution *Execution) createDependencyOrder() (dependencyOrder, error) {
	table := execution.populatePassEntryTable()
	if err := translatePassesToGraphEntries(table); err != nil {
		return dependencyOrder{}, err
	}
	entries := extractValues(table)
	return dependencyOrder{entries: entries}, nil
}

func extractValues(table map[Pass] *graphEntry) (values []*graphEntry) {
	for _, value := range table {
		values = append(values, value)
	}
	return values
}

func (execution *Execution) populatePassEntryTable() map[Pass] *graphEntry {
	entries := map[Pass] *graphEntry{}
	execution.traversePassDependencies(func(pass Pass, dependencies []Pass) {
		entries[pass] = &graphEntry{pass: pass}
	})
	return entries
}

func translatePassesToGraphEntries(entries map[Pass] *graphEntry) error {
	for _, entry := range entries {
		if err := translateDependencies(entry, entries); err != nil {
			return err
		}
	}
	return nil
}

func translateDependencies(
	entry *graphEntry, entries map[Pass] *graphEntry) error {

	for _, dependency := range entry.pass.Dependencies() {
		if dependencyEntry, exists := entries[dependency]; exists {
			entry.dependencies = append(entry.dependencies, dependencyEntry)
		} else {
			return fmt.Errorf("graph entry for %s doesn't exist", dependency)
		}
	}
	return nil
}

type dependencyVisitor func(pass Pass, dependencies []Pass)

func (execution *Execution) traversePassDependencies(visitor dependencyVisitor) {
	execution.traversePassDependenciesRecursive(execution.target, visitor)
}

func (execution *Execution) traversePassDependenciesRecursive(
	pass Pass, visitor dependencyVisitor) {

	visitor(pass, pass.Dependencies())
	for _, dependency := range pass.Dependencies() {
		visitor(dependency, dependency.Dependencies())
	}
}

type dependencyOrder struct {
	entries []*graphEntry
	output []Pass
}

type graphEntry struct {
	pass Pass
	dependencies []*graphEntry
	executed bool
}

func (order *dependencyOrder) insert(entry *graphEntry) {
	if !order.contains(entry.pass) {
		order.entries = append(order.entries, entry)
	}
}

func (order *dependencyOrder) contains(pass Pass) bool {
	for _, entry := range order.entries {
		if entry.pass == pass {
			return true
		}
	}
	return false
}

func (order *dependencyOrder) compute() (output []Pass, err error) {
	if order.output == nil || len(order.output) == 0 {
		err = order.computeOutput()
	}
	return order.output, err
}

func (order *dependencyOrder) computeOutput() error {
	order.output = make([]Pass, len(order.entries))
	for index := range order.output {
		if err := order.computeIndex(index); err != nil {
			return err
		}
	}
	return nil
}

func (order *dependencyOrder) computeIndex(index int) error {
	if next, ok := order.nextElementToExecute(); ok {
		next.executed = true
		order.output[index] = next.pass
		return nil
	}
	return errors.New("could not sort pass order")
}

func (order *dependencyOrder) nextElementToExecute() (*graphEntry, bool) {
	for _, entry := range order.entries {
		if !entry.executed && entry.canBeExecuted() {
			return entry, true
		}
	}
	return nil, false
}

func (entry *graphEntry) canBeExecuted() bool {
	for _, dependency  := range entry.dependencies {
		if !dependency.executed {
			return false
		}
	}
	return true
}
