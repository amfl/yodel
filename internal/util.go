package yodel

import (
	mapset "github.com/deckarep/golang-set"
)

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func AsSet(strList []string) mapset.Set {
	s := make([]interface{}, len(strList))
	for i, v := range strList {
		s[i] = v
	}
	return mapset.NewSetFromSlice(s)
}
