package yodel

import (
	mapset "github.com/deckarep/golang-set"
	"gopkg.in/yaml.v2"
)

// OutputYaml outputs a set of groups as yaml :)
func OutputYaml(groups mapset.Set, annotate func(string) string) (string, error) {
	// Convert to slice
	slice := groups.ToSlice()

	// Convert to yaml
	d, err := yaml.Marshal(&slice)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
