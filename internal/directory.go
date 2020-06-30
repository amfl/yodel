package yodel

import (
	mapset "github.com/deckarep/golang-set"
)

// GroupSet represents a set of groups, on which we can perform set operations
// (difference, intersection, etc).
type GroupSet mapset.Set

// Directory is an interface which represents a queryable directory service,
// such as LDAP.
type Directory interface {
	Search(lookup string) (GroupSet, error)
	Sync() error
}
