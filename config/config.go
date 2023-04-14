package config

import (
	"github.com/spf13/viper"
	"os"
	"study02-chat-service/log"
)

type Config struct {
	Redis  RedisConfig  `mapstructure:"redis"`
	Server ServerConfig `mapstructure:"server"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func LoadConfigOrExit() Config {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./cmd")
	viper.AddConfigPath("./cmd/config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Error reading config file, %s", err)
		os.Exit(1)
	}

	// logging config path
	log.Infof("Using config: %s", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&config); err != nil {
		log.Errorf("unable to decode into struct, %v", err)
		os.Exit(1)
	}

	return config
}
