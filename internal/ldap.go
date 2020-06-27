package yodel

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
)

func CNToGroupName(cn string) string {
	// Pretty hacky!
	// Transform "cn=adminUser,ou=groups,dc=org,dc=example" => "adminUser"
	i := strings.Index(cn, ",")
	return cn[3:i]
}

func GetGroupsFromLdap() []string {
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

	// Assume that the first matching user is always the right one
	entry := sr.Entries[0]
	groups := entry.GetAttributeValues(viper.GetString("ldap.group_attribute"))

	return Map(groups, CNToGroupName)
}
