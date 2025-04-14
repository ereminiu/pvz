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
	"github.com/ereminiu/pvz/internal/pkg/db"
	"github.com/ereminiu/pvz/internal/repository/pvz"
	"github.com/ereminiu/pvz/tests/pgcontainer"
	"github.com/ereminiu/pvz/tests/repository/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PVZSuite struct {
	suite.Suite

	pgContainer *postgres.PostgresContainer
	r           *pvz.Repository
	database    *db.Database
	ctx         context.Context
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

	fmt.Println(connString)

	database, err := db.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	suite.database = database
	suite.r = pvz.New(database)
}

func (suite *PVZSuite) TearDownSuite() {
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
	orders = append(orders, collections.DeliveredOrders...)
	orders = append(orders, collections.GivenOrders...)
	orders = append(orders, collections.RefundOrders...)

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

	testTable := []struct {
		name           string
		page           int
		limit          int
		orderBy        string
		expectedResult []*entities.Order
		expectedError  error
	}{
		{
			name:           "OK",
			page:           1,
			limit:          3,
			orderBy:        "",
			expectedResult: collections.RefundOrders,
			expectedError:  nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			res, err := suite.r.GetRefunds(suite.ctx, tc.page, tc.limit, tc.orderBy, map[string]string{})

			now := time.Now()
			for i := range res {
				res[i].CreatedAt = now
				res[i].UpdatedAt = now
				res[i].ExpireAt = now

				tc.expectedResult[i].CreatedAt = now
				tc.expectedResult[i].UpdatedAt = now
				tc.expectedResult[i].ExpireAt = now
			}

			assert.ErrorIs(t, err, tc.expectedError)
			assert.ElementsMatch(t, res, tc.expectedResult)
		})
	}
}

func (suite *PVZSuite) TestGetHistory() {
	t := suite.T()

	testTable := []struct {
		name           string
		page           int
		limit          int
		orderBy        string
		expectedResult []*entities.Order
		expectedError  error
	}{
		{
			name:           "OK",
			page:           1,
			limit:          3,
			orderBy:        "",
			expectedResult: collections.DeliveredOrders,
			expectedError:  nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			res, err := suite.r.GetHistory(suite.ctx, tc.page, tc.limit, tc.orderBy, map[string]string{})

			now := time.Now()
			for i := range res {
				res[i].CreatedAt = now
				res[i].UpdatedAt = now
				res[i].ExpireAt = now

				tc.expectedResult[i].CreatedAt = now
				tc.expectedResult[i].UpdatedAt = now
				tc.expectedResult[i].ExpireAt = now
			}

			assert.ErrorIs(t, err, tc.expectedError)
			assert.ElementsMatch(t, res, tc.expectedResult)
		})
	}
}

func (suite *PVZSuite) TestAboba() {
	t := suite.T()

	x := 12
	assert.Equal(t, 12, x)
}

func TestPVZSuite(t *testing.T) {
	suite.Run(t, new(PVZSuite))
}
