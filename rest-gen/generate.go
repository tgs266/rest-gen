package restgen

import (
	"os"
	"path/filepath"

	"github.com/tgs266/rest-gen/rest-gen/config"
	"github.com/tgs266/rest-gen/rest-gen/gen"
)

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func Generate(config config.Config) {
	inputYamlFiles, _ := WalkMatch(config.Definitions.InputDir, "*.yaml")
	inputYmlFiles, _ := WalkMatch(config.Definitions.InputDir, "*.yml")
	inputFiles := append(inputYamlFiles, inputYmlFiles...)
	if len(inputFiles) == 0 {
		panic("no input specs found")
	}
	os.MkdirAll(config.Definitions.OutputDir, os.ModePerm)
	for _, file := range inputFiles {
		gen.GenerateFromSpec(file, config.Definitions.OutputDir)
	}
}
