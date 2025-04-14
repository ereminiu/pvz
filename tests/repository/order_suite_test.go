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
	"github.com/ereminiu/pvz/internal/repository/order"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/ereminiu/pvz/tests/repository/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type OrderSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	r           *order.Repository
	database    *db.Database
	ctx         context.Context
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

	fmt.Println(connString)

	database, err := db.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	suite.database = database
	suite.r = order.New(database)
}

func (suite *OrderSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error during teardown: %v", err)
	}
}

func (suite *OrderSuite) SetupTest() {
	query := fmt.Sprintf(`TRUNCATE %s`, "orders")
	if _, err := suite.database.Exec(suite.ctx, query); err != nil {
		panic(err)
	}
}

func (suite *OrderSuite) TestAddOrder() {
	t := suite.T()

	testTable := []struct {
		name          string
		order         *entities.Order
		expectedError error
	}{
		{
			name: "Invalid user id",
			order: &entities.Order{
				UserID:   199,
				OrderID:  27,
				ExpireAt: time.Now().AddDate(0, 0, 3),
				Weight:   5,
				Price:    100,
				Packing:  "box",
				Extra:    true,
				Status:   "delivered",
			},
			expectedError: myerrors.ErrOrderAlreadyCreated,
		},
		{
			name:          "OK",
			order:         collections.Order1,
			expectedError: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := suite.r.AddOrder(suite.ctx, tc.order)

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func (suite *OrderSuite) TestRemoveOrder() {
	t := suite.T()

	testTable := []struct {
		name          string
		order         *entities.Order
		id            int
		expectedError error
	}{
		{
			name:          "OK",
			order:         collections.Order1,
			id:            89,
			expectedError: nil,
		},
		{
			name:          "OK",
			order:         collections.Order1,
			id:            89,
			expectedError: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			err := suite.r.AddOrder(suite.ctx, tc.order)

			assert.NoError(t, err)

			err = suite.r.RemoveOrder(suite.ctx, tc.id)

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func (suite *OrderSuite) TestAboba() {
	t := suite.T()

	x := 12
	assert.Equal(t, 12, x)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderSuite))
}
