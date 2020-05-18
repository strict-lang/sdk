package buildtool

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Configuration struct {
	PackageName  string                    `yaml:"package" json:"package"`
	Author       string                    `yaml:"author" json:"author"`
	Description  string                    `yaml:"description" json:"description"`
	Repositories []RepositoryConfiguration `yaml:"repositories" json:"repositories"`
}

type RepositoryConfiguration struct {
	Url  string `yaml:"url"`
	Name string `yaml:"name"`
}

func ReadConfiguration(file string) (Configuration, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return Configuration{}, fmt.Errorf("could not read build config: %v", err)
	}
  var configuration Configuration
	if err := yaml.Unmarshal(content, &configuration); err != nil {
		return Configuration{}, fmt.Errorf("could not parse build config: %v", err)
	}
	return configuration, nil
}
