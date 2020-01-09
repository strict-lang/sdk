package isolate

import (
	"sync"
)

type sharedFactory struct {
	lock sync.Mutex
	configurators []configurator
}

type configurator func(*Isolate)

var factory = &sharedFactory{}

func RegisterConfigurator(configurator configurator) {
	factory.lock.Lock()
	defer factory.lock.Unlock()

	factory.configurators = append(factory.configurators, configurator)
}

func (factory *sharedFactory) createIsolate() *Isolate {
	isolate := &Isolate{
		Properties: NewThreadLocalPropertyTable(),
	}
	for _, configurator := range factory.configurators {
		configurator(isolate)
	}
	return isolate
}

var cachedGlobalIsolate = struct {
	lock sync.Mutex
	isolate *Isolate
}{}

func SingleThreaded() *Isolate {
	cachedGlobalIsolate.lock.Lock()
	defer cachedGlobalIsolate.lock.Unlock()
	if cachedGlobalIsolate.isolate == nil {
		cachedGlobalIsolate.isolate = factory.createIsolate()
	}
	return cachedGlobalIsolate.isolate
}
