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
	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/handler/pvz"
	"github.com/ereminiu/pvz/internal/pkg/db"
	rep "github.com/ereminiu/pvz/internal/repository"
	uu "github.com/ereminiu/pvz/internal/usecases/pvz"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/ereminiu/pvz/tests/repository/collections"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

type PVZSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	database    *db.Database
	h           *pvz.Handler
	ctx         context.Context
}

func getPVZHandler(ctx context.Context, database *db.Database) *pvz.Handler {
	repository := rep.New(database)
	uc := uu.New(repository)
	pvzHandler := pvz.New(ctx, uc)

	return pvzHandler
}

func (suite *PVZSuite) SetupSuite() {
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
	pvzHandler := getPVZHandler(ctx, database)
	suite.h = pvzHandler
}

func (suite *PVZSuite) TearDownSuite() {
	suite.database.GetPool().Close()
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error during teardown: %v", err)
	}
}

func (suite *PVZSuite) SetupTest() {
	query := `TRUNCATE orders`
	if _, err := suite.database.Exec(suite.ctx, query); err != nil {
		panic(err)
	}

	orders := []*entities.Order{}
	orders = append(orders, collections.RefundOrders...)
	orders = append(orders, collections.DeliveredOrders...)

	query = `INSERT INTO 
			orders (user_id, id, expire_at, weight, price, packing, extra, status)
			values ($1, $2, $3, $4, $5, $6, $7, $8)`

	for _, order := range orders {
		_, err := suite.database.Exec(suite.ctx, query,
			order.UserID,
			order.OrderID,
			order.ExpireAt,
			order.Weight,
			order.Price,
			order.Packing,
			order.Extra,
			order.Status,
		)

		if err != nil {
			panic(err)
		}
	}
}

func (suite *PVZSuite) TestGetRefunds() {
	t := suite.T()

	type getRefundsInput struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		OrderBy string `json:"order_by"`
	}

	testTable := []struct {
		name               string
		input              *getRefundsInput
		expectedStatusCode int
	}{
		{
			name: "OK",
			input: &getRefundsInput{
				Page:    1,
				Limit:   4,
				OrderBy: "price",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid order_by",
			input: &getRefundsInput{
				Page:    1,
				Limit:   4,
				OrderBy: "ppprice",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /refund-list", suite.h.GetRefunds)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refund-list", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
		})
	}
}

func (suite *PVZSuite) TestGetHistory() {
	t := suite.T()

	type getHistoryInput struct {
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		OrderBy string `json:"order_by"`
	}

	testTable := []struct {
		name                 string
		input                *getHistoryInput
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: &getHistoryInput{
				Page:    1,
				Limit:   4,
				OrderBy: "price",
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid order_by",
			input: &getHistoryInput{
				Page:    1,
				Limit:   4,
				OrderBy: "ppprice",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /history", suite.h.GetHistory)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/history", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
		})
	}
}

func TestPVZSuite(t *testing.T) {
	suite.Run(t, new(PVZSuite))
}
