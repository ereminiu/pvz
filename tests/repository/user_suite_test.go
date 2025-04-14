//go:build integration
// +build integration

package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ereminiu/pvz/internal/entities"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pkg/db"
	"github.com/ereminiu/pvz/internal/repository/user"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/ereminiu/pvz/tests/repository/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type UserSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	r           *user.Repository
	database    *db.Database
	ctx         context.Context
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

	fmt.Println(connString)

	database, err := db.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	suite.database = database
	suite.r = user.New(database)
}

func (suite *UserSuite) TearDownSuite() {
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

	testTable := []struct {
		name          string
		userID        int
		orderIDs      []int
		expectedError error
	}{
		{
			name:          "OK",
			userID:        26,
			orderIDs:      []int{94, 93},
			expectedError: nil,
		},
		{
			name:          "Invalid userID",
			userID:        16,
			orderIDs:      []int{29},
			expectedError: myerrors.ErrInvalidOrderInput,
		},
		{
			name:          "Empty orderIDs",
			userID:        26,
			orderIDs:      []int{},
			expectedError: myerrors.ErrInvalidOrderInput,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := suite.r.RefundOrders(suite.ctx, tc.userID, tc.orderIDs)

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func (suite *UserSuite) TestReturnOrders() {
	t := suite.T()

	testTable := []struct {
		name          string
		userID        int
		orderIDs      []int
		expectedError error
	}{
		{
			name:          "OK",
			userID:        19,
			orderIDs:      []int{89, 86, 27},
			expectedError: nil,
		},
		{
			name:          "Invalid userID",
			userID:        199,
			orderIDs:      []int{11, 22},
			expectedError: myerrors.ErrInvalidOrderInput,
		},
		{
			name:          "Empty orderIDs",
			userID:        19,
			orderIDs:      []int{},
			expectedError: myerrors.ErrInvalidOrderInput,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := suite.r.ReturnOrders(suite.ctx, tc.userID, tc.orderIDs)

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func (suite *UserSuite) TestGetList() {
	t := suite.T()

	testTable := []struct {
		name           string
		userID         int
		lastN          int
		located        bool
		expectedResult []*entities.Order
		expectedError  error
	}{
		{
			name:    "OK",
			userID:  19,
			lastN:   -1,
			located: false,
			expectedResult: []*entities.Order{
				collections.Order1,
				collections.Order2,
				collections.Order3,
			},
			expectedError: nil,
		},
		{
			name:           "Invalid UserID",
			userID:         17,
			lastN:          -1,
			located:        false,
			expectedResult: nil,
			expectedError:  nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			res, err := suite.r.GetList(suite.ctx, tc.userID, tc.lastN, tc.located, map[string]string{})

			assert.Equal(t, len(res), len(tc.expectedResult))

			now := time.Now()
			for i := range res {
				res[i].CreatedAt = now
				res[i].UpdatedAt = now
				res[i].ExpireAt = now

				tc.expectedResult[i].CreatedAt = now
				tc.expectedResult[i].UpdatedAt = now
				tc.expectedResult[i].ExpireAt = now
			}

			assert.ElementsMatch(t, res, tc.expectedResult)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func (suite *UserSuite) TestAboba() {
	t := suite.T()

	x := 12

	assert.Equal(t, x, 12)
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}
