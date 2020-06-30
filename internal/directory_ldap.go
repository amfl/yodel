package yodel

import (
	"fmt"
	"log"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
)

type LdapConfig struct {
	HostURL        string
	BindDN         string
	BindPassword   string
	GroupAttribute string
	BaseDN         string // Base DN on which to search
	Filter         string // Filter which uniquely identifies the user
}

func CNToGroupName(cn string) string {
	// Pretty hacky!
	// Transform "cn=adminUser,ou=groups,dc=org,dc=example" => "adminUser"
	i := strings.Index(cn, ",")
	return cn[3:i]
}

func GenerateLdapConfig() LdapConfig {
	return LdapConfig{
		HostURL:        viper.GetString("ldap.host_url"),
		BindDN:         viper.GetString("ldap.bind_dn"),
		BindPassword:   viper.GetString("ldap.bind_password"),
		GroupAttribute: viper.GetString("ldap.group_attribute"),
		BaseDN:         viper.GetString("ldap.base_dn"),
		Filter:         viper.GetString("ldap.filter"),
	}
}

//////////////////////

type LdapDirectory struct {
	cache  GroupSet
	config LdapConfig
}

func NewLdapDirectory(config LdapConfig) *LdapDirectory {
	l := new(LdapDirectory)
	l.cache = mapset.NewSet()
	l.config = config
	return l
}

func (l LdapDirectory) Search(lookup string) (GroupSet, error) {
	log.Print(l.config.HostURL)
	ld, err := ldap.DialURL(l.config.HostURL)
	if err != nil {
		log.Fatal(err)
	}
	defer ld.Close()

	log.Print("Dialed")

	err = ld.Bind(l.config.BindDN, l.config.BindPassword)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Bound")

	// Define search
	filter := fmt.Sprintf(l.config.Filter, lookup)
	searchRequest := ldap.NewSearchRequest(
		l.config.BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,                                  // The filter to apply
		[]string{"cn", l.config.GroupAttribute}, // A list attributes to retrieve
		nil,
	)

	log.Print("Defined search")

	// Execute search
	sr, err := ld.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Searched")

	// Assume that the first matching user is always the right one
	entry := sr.Entries[0]
	groups := entry.GetAttributeValues(l.config.GroupAttribute)

	return AsSet(Map(groups, CNToGroupName)), nil
}

func (l LdapDirectory) Sync() error {
	return nil
}
