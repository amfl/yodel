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

func AsSet(str_list []string) mapset.Set {
	s := make([]interface{}, len(str_list))
	for i, v := range str_list {
		s[i] = v
	}
	return mapset.NewSetFromSlice(s)
}
