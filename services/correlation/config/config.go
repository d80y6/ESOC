package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	KafkaBrokers     []string `mapstructure:"KAFKA_BROKERS"`
	InputTopic       string   `mapstructure:"INPUT_TOPIC"`
	OutputTopic      string   `mapstructure:"OUTPUT_TOPIC"`
	ConsumerGroupID  string   `mapstructure:"CONSUMER_GROUP_ID"`
	LogLevel         string   `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("KAFKA_BROKERS", []string{"kafka:9092"})
	viper.SetDefault("INPUT_TOPIC", "normalized-events")
	viper.SetDefault("OUTPUT_TOPIC", "alerts")
	viper.SetDefault("CONSUMER_GROUP_ID", "correlation-engine")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
