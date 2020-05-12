package buildtool

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

type Configuration struct {
	PackageName  string                    `yaml:"package" json:"package"`
	Author       string                    `yaml:"author" json:"author"`
	Description  string                    `yaml:"description" json:"description"`
	Repositories []RepositoryConfiguration `yaml:"repositories" json:"repositories"`
}

type RepositoryConfiguration struct {
	Url  string `yaml:"url" json:"url"`
	Name string `yaml:"name" json:"name"`
}

func ReadConfiguration(file string) (Configuration, error) {
	function, err := createUnmarshalFunction(file)
	if err != nil	{
		return Configuration{}, err
	}
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return Configuration{}, err
	}
	configuration := &Configuration{}
	err = function(contents, configuration)
	return *configuration, err
}

type unmarshalFunction func(bytes []byte, target interface{}) error

func createUnmarshalFunction(fileName string) (unmarshalFunction, error) {
	switch path.Ext(fileName) {
	case "yml", "yaml":
		return yaml.Unmarshal, nil
	case "json":
		return json.Unmarshal, nil

	}
	return nil, errors.New("unsupported configuration type")
}