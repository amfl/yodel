package yodel

import (
	"io/ioutil"

	mapset "github.com/deckarep/golang-set"
	"gopkg.in/yaml.v2"
)

type Role struct {
	Name   string
	Groups []string
	Roles  []string
}

type Group struct {
	Name string
	Desc string
}

type GroupFile struct {
	SchemaVersion string `yaml:"schema_version"`
	Roles         []Role
	Groups        []Group
}

//////////////////////

type YamlDirectory struct {
	cache    GroupSet
	gf       GroupFile
	filepath string
}

func NewYamlDirectory(filepath string) *YamlDirectory {
	y := new(YamlDirectory)
	y.cache = mapset.NewSet()
	y.filepath = filepath
	return y
}

func (y *YamlDirectory) Search(lookup string) (GroupSet, error) {
	result := mapset.NewSet()
	for _, r := range y.gf.Roles {
		if r.Name == lookup {
			// Add all the groups in this role
			for _, groupName := range r.Groups {
				result.Add(groupName)
			}
			// Recursively add the groups in any subroles
			// TODO BUG CAUTION - No infinite recursion protection
			for _, subroleName := range r.Roles {
				subgroups, err := y.Search(subroleName)
				if err != nil {
					return result, err
				}
				result = result.Union(subgroups)
			}
		}
	}
	return result, nil
}

func (y *YamlDirectory) Sync() error {
	groupDB, err := ioutil.ReadFile(y.filepath)
	if err != nil {
		return err
	}

	gf := GroupFile{}
	err = yaml.Unmarshal(groupDB, &gf)
	if err != nil {
		return err
	}

	// Save the result
	y.gf = gf

	return nil
}

func (y YamlDirectory) AnnotationFunction(string) string {
	return "Temporary string!"
}
