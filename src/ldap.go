package main

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
	"log"
)

func GetGroupsFromLdap() {
	log.Print(viper.GetString("ldap.host_url"))
	l, err := ldap.DialURL(viper.GetString("ldap.host_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = l.Bind(viper.GetString("ldap.bind_dn"),
		viper.GetString("ldap.bind_password"))
	if err != nil {
		log.Fatal(err)
	}

	// Define search
	filter := fmt.Sprintf(viper.GetString("ldap.filter"), viper.GetString("ldap.user"))
	searchRequest := ldap.NewSearchRequest(
		viper.GetString("ldap.base_dn"), // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, // The filter to apply
		[]string{"cn", viper.GetString("ldap.group_attribute")}, // A list attributes to retrieve
		nil,
	)

	// Execute search
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range sr.Entries {
		// Print all attributes for this entry
		for _, val := range entry.GetAttributeValues(viper.GetString("ldap.group_attribute")) {
			fmt.Printf("%s: %s\n", entry.DN, val)
		}
	}
}
