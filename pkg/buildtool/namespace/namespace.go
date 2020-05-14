package namespace

type Namespace interface {
	Name() string
	QualifiedName() string
	Dependencies() []Dependency
	Entries() []Entry
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
	computed      bool
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

func (namespace *namespace) Dependencies() []Namespace {

	return nil
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
