package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	DatabaseDSN string `mapstructure:"DATABASE_DSN"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("HTTP_PORT", "8086")
	viper.SetDefault("DATABASE_DSN", "host=postgres user=omniguard password=omniguard dbname=omniguard port=5432 sslmode=disable")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
