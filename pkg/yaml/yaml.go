package yaml

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// Load unmarshals YAML data from the filesystem
func Load(path string, result interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, result); err != nil {
		return err
	}

	return nil
}
