package spec

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Read(path string) *Spec {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var spec *Spec
	err = yaml.UnmarshalStrict(file, &spec)
	if err != nil {
		panic(err)
	}
	return spec
}
