package yodel

import (
	"io/ioutil"

	mapset "github.com/deckarep/golang-set"
	"gopkg.in/yaml.v2"
)

// YamlRole is used when deserializing group database files.
type YamlRole struct {
	Name   string
	Groups []string
	Roles  []string
}

// YamlGroup is used when deserializing group database files.
type YamlGroup struct {
	Name string
	Desc string
}

// YamlGroupFile is the in-memory representation of a group database file.
type YamlGroupFile struct {
	SchemaVersion string `yaml:"schema_version"`
	Roles         []YamlRole
	Groups        []YamlGroup
}

//////////////////////

// YamlDirectory exposes the group database file as a query directory service.
type YamlDirectory struct {
	cache    GroupSet
	gf       YamlGroupFile
	filepath string
}

// NewYamlDirectory acts as a constructor for YamlDirectory
func NewYamlDirectory(filepath string) *YamlDirectory {
	y := new(YamlDirectory)
	y.cache = mapset.NewSet()
	y.filepath = filepath
	return y
}

// Search performs a search on the in-memory group database by finding roles
// which match the given `lookup` string.
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

// Sync gets this directory service ready to issue searches. In the case of
// YAML, this method reads the database file into memory.
func (y *YamlDirectory) Sync() error {
	groupDB, err := ioutil.ReadFile(y.filepath)
	if err != nil {
		return err
	}

	gf := YamlGroupFile{}
	err = yaml.Unmarshal(groupDB, &gf)
	if err != nil {
		return err
	}

	// Save the result
	y.gf = gf

	return nil
}

// AnnotationFunction can be used to associate groups with descriptions
// present in the group database file. Useful to make a list of groups more
// human readable in program output.
func (y YamlDirectory) AnnotationFunction(string) string {
	// TODO
	return "UNIMPLEMENTED"
}
