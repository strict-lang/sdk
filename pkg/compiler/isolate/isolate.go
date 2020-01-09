package isolate

// Isolate is a per thread data structure that contains instances which can not
// be shared across different threads. It is used to run some processing and
// analysis concurrent, without extra synchronization.
type Isolate struct {
	Properties *ThreadLocalPropertyTable
}

func RegisterConfigurator(configurator func(*Isolate)) {
}
