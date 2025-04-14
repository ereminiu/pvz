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
	"github.com/ereminiu/pvz/internal/handler/user"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/db"
	rep "github.com/ereminiu/pvz/internal/repository"
	uu "github.com/ereminiu/pvz/internal/usecases/user"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/ereminiu/pvz/tests/repository/collections"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gotest.tools/v3/assert"
)

type UserSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	database    *db.Database
	h           *user.Handler
	ctx         context.Context
}

func getUserHandler(ctx context.Context, database *db.Database) *user.Handler {
	repository := rep.New(database)
	uc := uu.New(repository)
	userHandler := user.New(ctx, uc)

	return userHandler
}

func (suite *UserSuite) SetupSuite() {
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
	userHandler := getUserHandler(ctx, database)
	suite.h = userHandler
}

func (suite *UserSuite) TearDownSuite() {
	suite.database.GetPool().Close()
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error during teardown: %v", err)
	}
}

func (suite *UserSuite) SetupTest() {
	query := `TRUNCATE orders`
	if _, err := suite.database.Exec(suite.ctx, query); err != nil {
		panic(err)
	}

	orders := []*entities.Order{}
	orders = append(orders, collections.DeliveredOrders...)
	orders = append(orders, collections.GivenOrders...)

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

func (suite *UserSuite) TestRefundOrders() {
	t := suite.T()

	type inputRefundOrder struct {
		UserID   int   `json:"user_id"`
		OrderIDs []int `json:"order_ids"`
	}

	testTable := []struct {
		name                 string
		input                *inputRefundOrder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: &inputRefundOrder{
				UserID:   26,
				OrderIDs: []int{94, 93},
			},
			expectedStatusCode:   http.StatusResetContent,
			expectedResponseBody: "orders have been refunded",
		},
		{
			name: "Invalid order input",
			input: &inputRefundOrder{
				UserID:   62,
				OrderIDs: []int{49, 39},
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrInvalidOrderInput.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /refund", suite.h.RefundOrders)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refund", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func (suite *UserSuite) TestReturnOrders() {
	t := suite.T()

	type inputRefundOrder struct {
		UserID   int   `json:"user_id"`
		OrderIDs []int `json:"order_ids"`
	}

	testTable := []struct {
		name                 string
		input                *inputRefundOrder
		expectedStatusCode   int
		expectedError        error
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: &inputRefundOrder{
				UserID:   19,
				OrderIDs: []int{89, 86, 27},
			},
			expectedStatusCode:   http.StatusResetContent,
			expectedResponseBody: "orders have been returned",
		},
		{
			name: "Invalid order input",
			input: &inputRefundOrder{
				UserID:   62,
				OrderIDs: []int{49, 39},
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrInvalidOrderInput.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /return", suite.h.ReturnOrders)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/return", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func (suite *UserSuite) TestGetList() {
	t := suite.T()

	type inputGetList struct {
		UserID  int  `json:"user_id"`
		LastN   int  `json:"last_n"`
		Located bool `json:"located"`
	}

	testTable := []struct {
		name               string
		input              *inputGetList
		expectedStatusCode int
	}{
		{
			name: "OK",
			input: &inputGetList{
				UserID:  19,
				LastN:   -1,
				Located: true,
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			router := http.NewServeMux()
			router.HandleFunc("POST /list", suite.h.GetList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/list", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tc.expectedStatusCode)
		})
	}
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
