// Based on go-ircevent examples
// https://github.com/thoj/go-ircevent/tree/master/examples

package main

import (
	"fmt"
	"log"
	"os"
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

func main() {
	readConfigs()

	groups := yodel.GetGroupsFromLdap()

	// DEBUGGING - Print all attributes for this entry
	for _, group := range groups {
		log.Print(group)
	}

	group_db := yodel.GetGroupsFromYaml(viper.GetString("groups.file"))

	yodel.Crunch("Wizard", groups, group_db)

	log.Print("All done")
}
