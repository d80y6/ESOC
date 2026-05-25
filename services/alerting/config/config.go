package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	KafkaBrokers    []string `mapstructure:"KAFKA_BROKERS"`
	AlertsTopic     string   `mapstructure:"ALERTS_TOPIC"`
	ConsumerGroupID string   `mapstructure:"CONSUMER_GROUP_ID"`
	DatabaseDSN     string   `mapstructure:"DATABASE_DSN"`
	HTTPPort        string   `mapstructure:"HTTP_PORT"`
	OIDCProviderURL string   `mapstructure:"OIDC_PROVIDER_URL"`
	OIDCAudience    string   `mapstructure:"OIDC_AUDIENCE"`
	LogLevel        string   `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("KAFKA_BROKERS", []string{"kafka:9092"})
	viper.SetDefault("ALERTS_TOPIC", "alerts")
	viper.SetDefault("CONSUMER_GROUP_ID", "alerting-service")
	viper.SetDefault("DATABASE_DSN", "host=postgres user=omniguard password=omniguard dbname=omniguard port=5432 sslmode=disable")
	viper.SetDefault("HTTP_PORT", "8085")
	viper.SetDefault("OIDC_PROVIDER_URL", "http://keycloak:8080/realms/omniguard")
	viper.SetDefault("OIDC_AUDIENCE", "omniguard-backend")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
