package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort     string `mapstructure:"GRPC_PORT"`
	HTTPPort     string `mapstructure:"HTTP_PORT"`
	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
	RawLogsTopic string `mapstructure:"RAW_LOGS_TOPIC"`
	LogLevel     string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("GRPC_PORT", "8082")
	viper.SetDefault("HTTP_PORT", "8080")
	viper.SetDefault("KAFKA_BROKERS", []string{"kafka:9092"})
	viper.SetDefault("RAW_LOGS_TOPIC", "raw-logs")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
