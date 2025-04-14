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
	"github.com/ereminiu/pvz/internal/handler/order"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/db"
	rep "github.com/ereminiu/pvz/internal/repository"
	uu "github.com/ereminiu/pvz/internal/usecases/order"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/ereminiu/pvz/tests/repository/collections"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

type OrderSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	database    *db.Database
	h           *order.Handler
	ctx         context.Context
}

func getOrderHandler(ctx context.Context, database *db.Database) *order.Handler {
	repository := rep.New(database)
	uc := uu.New(repository)
	orderHandler := order.New(ctx, uc)

	return orderHandler
}

func (suite *OrderSuite) SetupSuite() {
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
	orderHandler := getOrderHandler(ctx, database)
	suite.h = orderHandler
}

func (suite *OrderSuite) TearDownSuite() {
	suite.database.GetPool().Close()
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error during teardown: %v", err)
	}
}

func (suite *OrderSuite) SetupTest() {
	query := `TRUNCATE orders`
	if _, err := suite.database.Exec(suite.ctx, query); err != nil {
		panic(err)
	}

	orders := []*entities.Order{}
	orders = append(orders, collections.Order2)

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

func (suite *OrderSuite) TestAdd() {
	t := suite.T()

	testTable := []struct {
		name                 string
		input                *entities.Order
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: &entities.Order{
				UserID:      19,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       40,
				Packing:     "box",
				Extra:       true,
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "order has been added",
		},
		{
			name: "Invalid UserID",
			input: &entities.Order{
				UserID:      91,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       40,
				Packing:     "box",
				Extra:       true,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrOrderAlreadyCreated.Error() + "\n",
		},
		{
			name: "Invalid packing type",
			input: &entities.Order{
				UserID:      19,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       40,
				Packing:     "xbox",
				Extra:       true,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrInvalidOrderPackingType.Error() + "\n",
		},
		{
			name: "Order already created",
			input: &entities.Order{
				UserID:      19,
				OrderID:     86,
				ExpireAfter: 2,
				Weight:      5,
				Price:       40,
				Packing:     "film",
				Extra:       false,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrOrderAlreadyCreated.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /add", suite.h.Add)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/add", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func (suite *OrderSuite) TestRemove() {
	t := suite.T()

	type removeOrderInput struct {
		ID int `json:"order_id"`
	}

	testTable := []struct {
		name                 string
		input                removeOrderInput
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: removeOrderInput{
				ID: 86,
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "order has been removed",
		},
		{
			name: "Order already removed",
			input: removeOrderInput{
				ID: 11,
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrOrderAlreadyRemoved.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("DELETE /remove", suite.h.Remove)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/remove", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func (suite *OrderSuite) TestAboba() {
	t := suite.T()

	x := 12

	assert.Equal(t, x, 12)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
