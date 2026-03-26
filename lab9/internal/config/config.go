package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	MongoDB MongoDB `mapstructure:"mongodb"`
}

type MongoDB struct {
	URI            string `mapstructure:"uri"`
	Database       string `mapstructure:"database"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

func Load(path string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(path)

	v.SetDefault("mongodb.timeout_seconds", 10)

	v.SetEnvPrefix("MONGODB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config file %s: %w", path, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) MongoTimeout() time.Duration {
	return time.Duration(c.MongoDB.TimeoutSeconds) * time.Second
}

func (c *Config) validate() error {
	c.MongoDB.URI = strings.TrimSpace(c.MongoDB.URI)
	c.MongoDB.Database = strings.TrimSpace(c.MongoDB.Database)

	if c.MongoDB.URI == "" {
		return fmt.Errorf("config validation: mongodb.uri is required")
	}

	if c.MongoDB.Database == "" {
		return fmt.Errorf("config validation: mongodb.database is required")
	}

	if c.MongoDB.TimeoutSeconds < 1 {
		return fmt.Errorf("config validation: mongodb.timeout_seconds must be at least 1")
	}

	return nil
}
