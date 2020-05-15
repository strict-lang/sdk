package namespace

type Namespace interface {
	Name() string
	QualifiedName() string
	Entries() []Entry
	MarkAsCompiled()
	IsCompiled() bool
}

type Dependency interface {
	Load()
}

type Entry interface {
	FileName() string
	IsDirectory() bool
}

type namespace struct {
	name string
	qualifiedName string
	entries       []Entry
	compiled bool
}

func (namespace *namespace) MarkAsCompiled() {
	namespace.compiled = true
}

func (namespace *namespace) IsCompiled() bool {
	return namespace.compiled
}

func (namespace *namespace) Name() string {
	return namespace.name
}

func (namespace *namespace) QualifiedName() string {
	return namespace.qualifiedName
}

func (namespace *namespace) Entries() []Entry {
	return namespace.entries
}

type entry struct {
	fileName string
	directory bool
}

func (entry *entry) FileName() string {
	return entry.fileName
}

func (entry *entry) IsDirectory() bool {
	return entry.directory
}


func NewRoot(directory string) Namespace {
	return nil
}
