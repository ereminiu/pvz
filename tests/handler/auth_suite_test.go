//go:build integration
// +build integration

package handler

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/ereminiu/pvz/internal/handler/auth"
	"github.com/ereminiu/pvz/internal/pkg/db"
	rep "github.com/ereminiu/pvz/internal/repository"
	uu "github.com/ereminiu/pvz/internal/usecases/auth"
	"github.com/ereminiu/pvz/internal/usecases/auth/hashgen"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

type AuthSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	database    *db.Database
	h           *auth.Handler
	ctx         context.Context
}

func getAuthHandler(ctx context.Context, database *db.Database) *auth.Handler {
	repository := rep.New(database)
	uc := uu.New(repository)
	authHandler := auth.New(ctx, uc)

	return authHandler
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

	database, err := db.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	suite.database = database
	authHandler := getAuthHandler(ctx, database)
	suite.h = authHandler
}

func (suite *AuthSuite) TearDownSuite() {
	suite.database.GetPool().Close()
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error during teardown: %v", err)
	}
}

func (suite *AuthSuite) SetupTest() {
	query := `TRUNCATE admins`
	if _, err := suite.database.Exec(suite.ctx, query); err != nil {
		panic(err)
	}

	query = `INSERT INTO 
			admins (username, password)
			VALUES ($1, $2)`

	username := "Alina Rin"
	password := hashgen.GenerateHash("qwerty128")

	if _, err := suite.database.Exec(suite.ctx, query, username, password); err != nil {
		panic(err)
	}
}

func (suite *AuthSuite) TestSignIn() {
	t := suite.T()

	type singInInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	testTable := []struct {
		name               string
		input              singInInput
		expectedStatusCode int
	}{
		{
			name: "OK",
			input: singInInput{
				Username: "Alina Rin",
				Password: "qwerty128",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid password",
			input: singInInput{
				Username: "Alina Rin",
				Password: "qwerty256",
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /sign-in", suite.h.SignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
		})
	}
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
