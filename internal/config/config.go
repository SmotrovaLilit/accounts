package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Service struct {
	Host string `mapstructure:"host" validate:"required"`
	Port int    `mapstructure:"port" validate:"required,min=1024,max=65535"`
}

type Prometheus struct {
	Port int `mapstructure:"port" validate:"required,min=1024,max=65535"`
}

type Config struct {
	Service     Service    `mapstructure:"service" validate:"required"`
	Prometheus  Prometheus `mapstructure:"prometheus" validate:"required"`
	DatabaseDSN string     `mapstructure:"database_dsn" validate:"required,url"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MYAPP") // TODO use the same app name like in traces and logs

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	return &config, nil
}
