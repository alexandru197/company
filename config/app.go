package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Server struct {
		Port int
	}
	JWT struct {
		Secret string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Kafka struct {
		Brokers []string
		Topic   string
	}
}

// InitConfig loads the configuration from config.yaml.
func InitConfig() Config {
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	return config
}
