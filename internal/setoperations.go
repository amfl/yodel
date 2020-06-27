package yodel

import (
	"fmt"

	mapset "github.com/deckarep/golang-set"
)

// func GroupLookupTable(group_db GroupFile) map[string]Group {
//     m := make(map[string]Group)

//     // Populate lookup table with all groups
//     for _, g := range group_db.Groups {
//         m[g.Name] = g
//     }

//     return m
// }

// func GroupNamesToGroup(group_names []string, lookup_name map[string]Group) []Group {
//     // TODO
//     return []Group{}
// }

func GetDesiredGroups(role_name string, group_db GroupFile) []string {
	group_names := []string{}
	for _, r := range group_db.Roles {
		if r.Name == role_name {
			// Add all the groups in this role
			for _, group_name := range r.Groups {
				group_names = append(group_names, group_name)
			}
			// Recursively add the groups in any subroles
			// TODO BUG CAUTION - No infinite recursion protection
			for _, subrole_name := range r.Roles {
				group_names = append(group_names, GetDesiredGroups(subrole_name, group_db)...)
			}
		}
	}
	return group_names
}

func AsSet(str_list []string) mapset.Set {
	s := make([]interface{}, len(str_list))
	for i, v := range str_list {
		s[i] = v
	}
	return mapset.NewSetFromSlice(s)
}

func Crunch(desired_role string, my_group_names []string, group_db GroupFile) {
	// lookup_table := GroupLookupTable(group_db)

	desired_group_names := GetDesiredGroups(desired_role, group_db)

	all_group_names := []string{}
	for _, g := range group_db.Groups {
		all_group_names = append(all_group_names, g.Name)
	}

	// Do set operations!

	fmt.Println(desired_group_names)
	needed_group_names := AsSet(desired_group_names).Difference(AsSet(my_group_names))

	fmt.Println("Your current groups:")
	fmt.Println(my_group_names)
	fmt.Println("Your role requires:")
	fmt.Println(desired_group_names)
	fmt.Println("Needed groups:")
	fmt.Println(needed_group_names)
}
