package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort     string `mapstructure:"GRPC_PORT"`
	HTTPPort     string `mapstructure:"HTTP_PORT"`
	KeycloakURL  string `mapstructure:"KEYCLOAK_URL"`
	KeycloakRealm string `mapstructure:"KEYCLOAK_REALM"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("GRPC_PORT", "8081")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("KEYCLOAK_URL", "http://keycloak:8080")
	viper.SetDefault("KEYCLOAK_REALM", "omniguard")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
