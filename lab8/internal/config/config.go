package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultPort = 61613

type Config struct {
	Broker       Broker       `json:"broker"`
	Destinations Destinations `json:"destinations"`
}

type Broker struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	HostName string `json:"host_name"`
}

type Destinations struct {
	Queue string `json:"queue"`
	Topic string `json:"topic"`
}

func Load(path string) (Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Config{}, fmt.Errorf("resolve config path: %w", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return Config{}, fmt.Errorf("open config file %s: %w", absPath, err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config file %s: %w", absPath, err)
	}

	if err := cfg.normalize(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Address() string {
	return net.JoinHostPort(c.Broker.Host, strconv.Itoa(c.Broker.Port))
}

func (c Config) QueueDestination() string {
	return destinationWithPrefix(c.Destinations.Queue, "/queue/")
}

func (c Config) TopicDestination() string {
	return destinationWithPrefix(c.Destinations.Topic, "/topic/")
}

func (c *Config) normalize() error {
	c.Broker.Host = strings.TrimSpace(c.Broker.Host)
	c.Broker.Username = strings.TrimSpace(c.Broker.Username)
	c.Broker.Password = strings.TrimSpace(c.Broker.Password)
	c.Broker.HostName = strings.TrimSpace(c.Broker.HostName)
	c.Destinations.Queue = strings.TrimSpace(c.Destinations.Queue)
	c.Destinations.Topic = strings.TrimSpace(c.Destinations.Topic)

	if c.Broker.Host == "" {
		return fmt.Errorf("config validation: broker.host is required")
	}

	if c.Broker.Port == 0 {
		c.Broker.Port = defaultPort
	}

	if c.Broker.Port < 1 || c.Broker.Port > 65535 {
		return fmt.Errorf("config validation: broker.port must be between 1 and 65535")
	}

	if c.Destinations.Queue == "" && c.Destinations.Topic == "" {
		return fmt.Errorf("config validation: at least one of destinations.queue or destinations.topic is required")
	}

	return nil
}

func destinationWithPrefix(name, prefix string) string {
	if name == "" {
		return ""
	}

	if strings.HasPrefix(name, "/") {
		return name
	}

	return prefix + name
}
