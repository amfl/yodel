package yodel

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)

// LdapConfig represents all the configuration required to connect to an LDAP
// server and run queries.
type LdapConfig struct {
	HostURL        string
	BindDN         string
	BindPassword   string
	GroupAttribute string
	BaseDN         string // Base DN on which to search
	Filter         string // Filter which uniquely identifies the user
}

// CNToGroupName will transform a full CN string
// ("cn=adminUser,ou=groups,dc=org,dc=example") to a group name ("adminUser").
func CNToGroupName(cn string) string {
	// Pretty hacky!
	i := strings.Index(cn, ",")
	return cn[3:i]
}

// GenerateLdapConfig generates an LDAP config object from external config
// files or environment variables.
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

// LdapDirectory represents a queryable LDAP directory service,
type LdapDirectory struct {
	cache  GroupSet
	config LdapConfig
}

// NewLdapDirectory acts as a constructor for LdapDirectory
func NewLdapDirectory(config LdapConfig) *LdapDirectory {
	l := new(LdapDirectory)
	l.cache = mapset.NewSet()
	l.config = config
	return l
}

func getTLSConfig() *tls.Config {
	insecure := viper.GetBool("security.ignore_tls")

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	localCertFile := viper.GetString("security.certificate_authority")
	certs, err := ioutil.ReadFile(localCertFile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: insecure,
		RootCAs:            rootCAs,
	}

	return config
}

// Search performs a search against the configured LDAP server by
// substituting the `lookup` argument into the configured filter.
func (l LdapDirectory) Search(lookup string) (GroupSet, error) {
	log.Print(l.config.HostURL)
	dialOpt := ldap.DialWithTLSConfig(getTLSConfig())
	ld, err := ldap.DialURL(l.config.HostURL, dialOpt)
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

	resultSet := mapset.NewSet()
	for _, group := range groups {
		resultSet.Add(CNToGroupName(group))
	}
	return resultSet, nil
}

// Sync gets this directory service ready to issue searches. In the case of
// LDAP, this method does nothing.
func (l LdapDirectory) Sync() error {
	return nil
}
