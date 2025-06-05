package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment   string        `mapstructure:"ENVIRONMENT"`
	DatabaseUrl   string        `mapstructure:"DATABASE_URL"`
	Secretkey     string        `mapstructure:"SECRETKEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func NewConfig(path, env string) (*Config, error) {
	viper.SetConfigName(".env." + env)
	viper.AddConfigPath(path)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
