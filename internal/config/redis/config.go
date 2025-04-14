package redis

import (
	"fmt"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
)

const (
	localPath = "/internal/config/redis/local.yaml"
	prodPath  = "/internal/config/redis/prod.yaml"
)

type Config struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	DB      int    `yaml:"db"`
	Enabled bool   `yaml:"enabled"`
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

// TODO: допилить
func (c *Config) GetOpts() *redis.Options {
	return &redis.Options{
		Addr: c.GetAddress(),
		DB:   c.DB,
	}
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
