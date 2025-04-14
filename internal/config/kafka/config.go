package kafka

import (
	"fmt"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	localPath = "/internal/config/kafka/local.yaml"
	prodPath  = "/internal/config/kafka/prod.yaml"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func New() (Config, error) {
	var cfg Config
	path := getPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("file not exists: %w", err)
	}

	err := cleanenv.ReadConfig(path, &cfg)

	return cfg, err
}

func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getPath() string {
	var path string

	wd, _ := os.Getwd()
	wd = strings.TrimSuffix(wd, "/cmd")

	env := getEnv()
	switch env {
	case "local":
		path = localPath

	case "prod":
		path = prodPath

	default:
		panic("AAAAAAAAAAAAAAAa")
	}

	return wd + path
}

func getEnv() string {
	env := os.Getenv("ENV")
	if env == "" {
		return "local"
	}

	return env
}
