package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Definitions Defintions `yaml:"definitions"`
}

type Defintions struct {
	InputDir  string `yaml:"input-dir"`
	OutputDir string `yaml:"output-dir"`
}

func Read(path string) Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}
	return config
}
