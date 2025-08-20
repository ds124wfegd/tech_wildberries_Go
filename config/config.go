// Ininicializing common application configuration
package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct { // creating configuration structure
	Server   ServerConfig
	Postgres PostgresConfig
}

type ServerConfig struct { // creating server's config
	AppVersion   string `json:"appVersion"`
	Host         string `json:"host" validate:"required"`
	Port         string `json:"port" validate:"required"`
	Timeout      time.Duration
	Idle_timeout time.Duration
	Env          string `json:"environment"`
}

type PostgresConfig struct { // creating database's config
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"-"`
	DBName   string `json:"DBName"`
	SSLMode  string `json:"sslMode"`
	PgDriver string `json:"pgDriver"`
}

func LoadConfig() (*viper.Viper, error) { // the function of uploading a config file

	viperInstance := viper.New() // creating a viper instance

	viperInstance.AddConfigPath("./config") // setting the path to config's file
	viperInstance.SetConfigName("config")   // setting the name of the config file
	viperInstance.SetConfigType("yaml")     // setting the type of the config

	err := viperInstance.ReadInConfig() // the process of reading the config file

	if err != nil { // errors handling
		return nil, err
	}
	return viperInstance, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) { // the function of parsing a congig file

	var c Config

	err := v.Unmarshal(&c) // config unmarshaling
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err) // errors handling
		return nil, err
	}
	return &c, nil
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
