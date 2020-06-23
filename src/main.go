// Based on go-ircevent examples
// https://github.com/thoj/go-ircevent/tree/master/examples

package main

import (
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func main() {
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

	fmt.Println("All done")
}
