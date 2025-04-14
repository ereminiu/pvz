//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"log"
	"testing"

	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/db"
	"github.com/ereminiu/pvz/internal/repository/auth"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type AuthSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	r           *auth.Repository
	database    *db.Database
	ctx         context.Context
}

func (suite *AuthSuite) SetupSuite() {
	ctx := context.Background()
	suite.ctx = ctx

	container, err := pgcontainer.New(ctx)
	if err != nil {
		panic(err)
	}

	suite.pgContainer = container
	connString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	fmt.Println(connString)

	database, err := db.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	suite.database = database
	suite.r = auth.New(database)
}

func (suite *AuthSuite) TearDownSuite() {
	suite.database.GetPool().Close()
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error during teardown: %v", err)
	}
}

func (suite *AuthSuite) SetupTest() {
	query := fmt.Sprintf(`TRUNCATE %s`, "admins")
	if _, err := suite.database.Exec(suite.ctx, query); err != nil {
		panic(err)
	}

	query = `INSERT INTO 
			admins (username, password) 
			VALUES ($1, $2)`

	username := "Alina Rin"
	password := "qwerty128"

	if _, err := suite.database.Exec(suite.ctx, query, username, password); err != nil {
		panic(err)
	}
}

func (suite *AuthSuite) TestCreateUser() {
	t := suite.T()

	testTable := []struct {
		name          string
		username      string
		password      string
		expectedError error
	}{
		{
			name:          "OK: successful login",
			username:      "Alina Rin",
			password:      "qwerty128",
			expectedError: nil,
		},
		{
			name:          "OK: user created",
			username:      "Valera",
			password:      "qwery128",
			expectedError: nil,
		},
		{
			name:          "Invalid password",
			username:      "Alina Rin",
			password:      "qwerty129",
			expectedError: myerrors.ErrIncorrectPassword,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := suite.r.CreateUser(suite.ctx, tc.username, tc.password)

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
