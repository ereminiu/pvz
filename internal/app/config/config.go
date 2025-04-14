package config

import (
	"log"

	"github.com/ereminiu/pvz/internal/config/application"
	"github.com/ereminiu/pvz/internal/config/auditlog"
	"github.com/ereminiu/pvz/internal/config/kafka"
	"github.com/ereminiu/pvz/internal/config/postgres"
	"github.com/ereminiu/pvz/internal/config/redis"
)

type Config struct {
	Postgres postgres.Config
	Redis    redis.Config
	App      application.Config
	Audit    auditlog.Config
	Kafka    kafka.Config
}

func loadDBConfig() postgres.Config {
	config, err := postgres.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadRedisConfig() redis.Config {
	config, err := redis.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadAppConfig() application.Config {
	config, err := application.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadAuditLogConfig() auditlog.Config {
	config, err := auditlog.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func loadKafkaConfig() kafka.Config {
	config, err := kafka.New()
	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func Mustload() Config {
	return Config{
		Postgres: loadDBConfig(),
		Redis:    loadRedisConfig(),
		App:      loadAppConfig(),
		Audit:    loadAuditLogConfig(),
		Kafka:    loadKafkaConfig(),
	}
}
