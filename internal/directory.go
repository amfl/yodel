package yodel

import (
	mapset "github.com/deckarep/golang-set"
)

type groups mapset.Set

type Directory interface {
	Search(lookup string) (groups, error)
	Sync() error
}
