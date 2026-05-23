package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort      string   `mapstructure:"HTTP_PORT"`
	OpenSearchURLs []string `mapstructure:"OPENSEARCH_URLS"`
	LogLevel      string   `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("HTTP_PORT", "8083")
	viper.SetDefault("OPENSEARCH_URLS", []string{"http://opensearch:9200"})
	viper.SetDefault("LOG_LEVEL", "info")

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
