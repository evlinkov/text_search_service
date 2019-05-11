package service

import (
	"github.com/spf13/viper"
	"log"
)

const (
	configurationFileName = "config.conf"
)

type Configuration struct {
	Address  string
	Port     string
	Filename string
}

func GetConfiguration() *Configuration {
	viper.SetConfigName(configurationFileName)
	viper.AddConfigPath("src/text_search_service/config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %v", err)
	}
	viper.WatchConfig()
	return unloadConfiguration()
}

func unloadConfiguration() *Configuration {
	return &Configuration{
		Address:  viper.GetString("server.address"),
		Port:     viper.GetString("server.port"),
		Filename: viper.GetString("filename"),
	}
}
