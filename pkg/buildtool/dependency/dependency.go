package dependency

import (
	"fmt"
)

type Dependency struct {
	name            string
	group           string
	version         string
	exportedModules []Module
	dependencies    []*Dependency
}

type Version struct {
	major int
	minor int
	patch int
}

func (version *Version) String() string {
	return fmt.Sprintf("%d:%d:%d", version.major, version.minor, version.patch)
}

func (dependency *Dependency) Id() string {
	return fmt.Sprintf(
		"%s:%s:%s",
		dependency.group,
		dependency.name,
		dependency.version)
}

type Module struct{}

const dependChainRootIndex = 0

// detectCirculars detects if the root dependency depends on any dependency that
// itself depends on the root, thus creating a circle.
func detectCirculars(root *Dependency) (chains []DependChain, found bool) {
	for _, dependency := range root.dependencies {
		chain, foundInPath := detectCircularsInPath(root, dependency)
		chains = append(chains, chain)
		if foundInPath {
			found = true
		}
	}
	return chains, found
}

func detectCircularsInPath(root *Dependency, target *Dependency) (DependChain, bool) {
	return DependChain{}, false
}
