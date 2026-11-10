package config

import (
	"log"
	"github.com/spf13/viper"
)

func InitConfig(fileName string) *viper.Viper {
	config := new()
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}

	return config
}
