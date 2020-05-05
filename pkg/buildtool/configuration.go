package buildtool

type Configuration struct {
	PackageName string `yaml:"package" json:"package"`
	Author string `yaml:"author" json:"author"`
	Description string `yaml:"description" json:"description"`
	Repositories []RepositoryConfiguration `yaml:"repositories" json:"repositories"`
}

type RepositoryConfiguration struct {
	Url string `yaml:"url" json:"url"`
	Name string `yaml:"name" json:"name"`
}


