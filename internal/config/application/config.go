package application

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	localPath = "/internal/config/application/local.yaml"
	prodPath  = "/internal/config/application/prod.yaml"
)

type Config struct {
	Host              string        `yaml:"host"`
	ServiceName       string        `yaml:"service_name"`
	TracerPort        int           `yaml:"tracer_port"`
	RestPort          int           `yaml:"rest_port"`
	GRPCPort          int           `yaml:"grpc_port"`
	MonitoringPort    int           `yaml:"monitoring_port"`
	Env               string        `yaml:"env"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	CacheUpateTimeout time.Duration `yaml:"cache_update"`
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
