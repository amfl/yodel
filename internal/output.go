package yodel

import (
	mapset "github.com/deckarep/golang-set"
	"gopkg.in/yaml.v2"
)

func OutputYaml(groups mapset.Set, f func(string) string) (string, error) {
	// Convert to slice
	slice := groups.ToSlice()

	// Convert to yaml
	d, err := yaml.Marshal(&slice)
	if err != nil {
		return "", err
	}

	return string(d) + f("temp"), nil
}
