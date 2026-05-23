package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	KafkaBrokers     []string `mapstructure:"KAFKA_BROKERS"`
	RawLogsTopic     string   `mapstructure:"RAW_LOGS_TOPIC"`
	NormalizedTopic  string   `mapstructure:"NORMALIZED_TOPIC"`
	ConsumerGroupID  string   `mapstructure:"CONSUMER_GROUP_ID"`
	LogLevel         string   `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("KAFKA_BROKERS", []string{"kafka:9092"})
	viper.SetDefault("RAW_LOGS_TOPIC", "raw-logs")
	viper.SetDefault("NORMALIZED_TOPIC", "normalized-events")
	viper.SetDefault("CONSUMER_GROUP_ID", "normalization-service")
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
