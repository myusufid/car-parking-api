package config

import (
	"fmt"
	_ "github.com/flashlabs/rootpath"
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	// Load configuration using Viper
	config.SetConfigType("env") // Use env format
	config.AutomaticEnv()
	config.SetConfigName(".env")
	config.AddConfigPath("./")
	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	return config
}
