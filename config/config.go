package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	RedisAddress     string
	KafkaBrokers     string
	ElasticSearchURL string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
}
