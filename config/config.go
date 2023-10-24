package config

import (
	"strings"

	"github.com/abstractbreazy/rate-limited-api/pkg/limiter"
	"github.com/abstractbreazy/rate-limited-api/pkg/redis"
	"github.com/abstractbreazy/rate-limited-api/pkg/server"

	"github.com/spf13/viper"
)

type Config struct {
	Server      *server.Config  `yaml:"server" mapstructure:"server"`
	RateLimiter *limiter.Config `yaml:"limiter" mapstructure:"limiter"`
	Redis       *redis.Config   `yaml:"redis" mapstructure:"redis"`
}

func New() (conf *Config) {
	conf = new(Config)
	conf.Server = server.NewConfig()
	conf.RateLimiter = limiter.NewConfig()
	conf.Redis = &redis.Config{}
	return
}

func (c *Config) Init() (*Config, error) {
	var v = viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.SetEnvPrefix("limiter")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	// Check env variables
	v.AutomaticEnv()
	err = v.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
