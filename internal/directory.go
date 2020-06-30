package yodel

import (
	mapset "github.com/deckarep/golang-set"
)

type GroupSet mapset.Set

type Directory interface {
	Search(lookup string) (GroupSet, error)
	Sync() error
}
