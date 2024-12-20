package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// InitializeConfig initializes the configuration using Viper
func InitializeConfig() {
	// Set the file name of the configurations file
	viper.SetConfigFile("/etc/ceebffgo.properties")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Set the file type
	viper.SetConfigType("json")

	// Read in the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file", err)
	}

	// Optionally watch for changes to the configuration file
	viper.WatchConfig()
}
