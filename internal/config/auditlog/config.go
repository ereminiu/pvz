package auditlog

import (
	"fmt"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	localPath = "/internal/config/auditlog/local.yaml"
	prodPath  = "/internal/config/auditlog/prod.yaml"
)

type Config struct {
	Adapter string `yaml:"adapter"`
	Filter  string `yaml:"filter"`
}

func New() (Config, error) {
	var cfg Config
	path := getPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("file not exists: %w", err)
	}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, err
	}

	if !validate(cfg) {
		return Config{}, fmt.Errorf("invalid config setup")
	}

	return cfg, nil
}

func validate(cfg Config) bool {
	if cfg.Adapter != "default" && cfg.Adapter != "db" {
		return false
	}

	if cfg.Adapter == "db" && cfg.Filter != "" {
		return false
	}

	return true
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
