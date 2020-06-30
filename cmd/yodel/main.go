package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/amfl/yodel/internal"
	"github.com/spf13/viper"
)

func readConfigs() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	val, present := os.LookupEnv("YODEL_CONFIG_PATH")
	if present {
		viper.AddConfigPath(val)
	}

	// Sensibly allow env var settings
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("yodel")
	// Allows "test.thing" to be set with "YODEL_TEST_THING"

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s\n", err))
	}
}

func printUsage() {
	fmt.Printf("Usage: %s LDAP_USERNAME ROLE_NAME\n",
		filepath.Base(os.Args[0]))
}

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		printUsage()
		os.Exit(99)
	}

	ldapUser := args[0]
	roleName := args[1]

	readConfigs()

	// Construct our two directory interfaces
	ldapConfig := yodel.GenerateLdapConfig()
	ldapDir := yodel.NewLdapDirectory(ldapConfig)
	yamlDir := yodel.NewYamlDirectory(viper.GetString("groups.file"))

	// SYNC
	err := ldapDir.Sync()
	if err != nil {
		panic(fmt.Errorf("Error syncing LDAP: %s\n", err))
	}
	err = yamlDir.Sync()
	if err != nil {
		panic(fmt.Errorf("Error syncing YAML: %s\n", err))
	}

	// SEARCH
	// Extract groups from directory interfaces
	ldapGroups, err := ldapDir.Search(ldapUser)
	if err != nil {
		panic(fmt.Errorf("Error searching LDAP: %s\n", err))
	}
	log.Println(ldapGroups)
	yamlGroups, err := yamlDir.Search(roleName)
	if err != nil {
		panic(fmt.Errorf("Error searching YAML: %s\n", err))
	}

	// Find the difference
	diff := yamlGroups.Difference(ldapGroups)

	// Annotation function from yaml
	output, err := yodel.OutputYaml(diff, yamlDir.AnnotationFunction)
	if err != nil {
		panic(fmt.Errorf("Error formatting output: %s\n", err))
	}
	fmt.Println(output)

	log.Print("All done")
}
