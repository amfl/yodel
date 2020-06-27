package yodel

import (
	"fmt"
	"io/ioutil"
	"log"

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

func GetGroupsFromYaml(filename string) GroupFile {
	group_db, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	g := GroupFile{}
	err = yaml.Unmarshal(group_db, &g)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- g:\n%v\n\n", g)

	return g
}
