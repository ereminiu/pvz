package pgcontainer

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func New(ctx context.Context) (*postgres.PostgresContainer, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	wd = strings.TrimSuffix(wd, "repository")
	wd = strings.TrimSuffix(wd, "handler")
	path := wd + filepath.Join("..", "init.sql")

	fmt.Println(path)

	pgContainer, err := postgres.Run(ctx,
		"postgres:latest",
		postgres.WithInitScripts(path),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return pgContainer, err
}
