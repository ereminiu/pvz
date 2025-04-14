//go:build unit
// +build unit

package order_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ereminiu/pvz/internal/entities"
	"github.com/ereminiu/pvz/internal/handler"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/usecases"
	mock_usecases "github.com/ereminiu/pvz/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_Add(t *testing.T) {
	type mockBehavior func(s *mock_usecases.MockOrderUsecases, ctx context.Context, order *entities.Order)

	testTable := []struct {
		name                 string
		inputBody            string
		inputOrder           *entities.Order
		call                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"user_id":19,"order_id":89,"expire_after":2,"weight":5,"price":100,"packing":"box","extra":true}`,
			inputOrder: &entities.Order{
				UserID:      19,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       100,
				Packing:     "box",
				Extra:       true,
			},
			call: func(s *mock_usecases.MockOrderUsecases, ctx context.Context, order *entities.Order) {
				s.EXPECT().AddOrder(ctx, order).Return(nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: "order has been added",
		},
		{
			name:      "Order already created",
			inputBody: `{"user_id":19,"order_id":89,"expire_after":2,"weight":5,"price":100,"packing":"bobobobox","extra":true}`,
			inputOrder: &entities.Order{
				UserID:      19,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       100,
				Packing:     "bobobobox",
				Extra:       true,
			},
			call: func(s *mock_usecases.MockOrderUsecases, ctx context.Context, order *entities.Order) {
				s.EXPECT().AddOrder(ctx, order).Return(myerrors.ErrInvalidOrderPackingType)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "invalid packing type\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			order := mock_usecases.NewMockOrderUsecases(c)
			tc.call(order, ctx, tc.inputOrder)

			uc := &usecases.Usecases{
				OrderUsecases: order,
			}
			handler := handler.New(ctx, uc)

			router := http.NewServeMux()
			router.HandleFunc("POST /add", handler.Add)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/add", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_Remove(t *testing.T) {
	type mockBehavior func(s *mock_usecases.MockOrderUsecases, ctx context.Context, orderID int)

	testTable := []struct {
		name                 string
		inputBody            string
		inputOrderID         int
		ctx                  context.Context
		call                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "OK",
			inputBody:    `{"order_id":89}`,
			inputOrderID: 89,
			ctx:          context.Background(),
			call: func(s *mock_usecases.MockOrderUsecases, ctx context.Context, orderID int) {
				s.EXPECT().RemoveOrder(ctx, orderID).Return(nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "order has been removed",
		},
		{
			name:         "Invalid order_id",
			inputBody:    `{"order_id":89}`,
			inputOrderID: 89,
			ctx:          context.Background(),
			call: func(s *mock_usecases.MockOrderUsecases, ctx context.Context, orderID int) {
				s.EXPECT().RemoveOrder(ctx, orderID).Return(myerrors.ErrInvalidUserInput)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrInvalidUserInput.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			order := mock_usecases.NewMockOrderUsecases(c)
			tc.call(order, tc.ctx, tc.inputOrderID)

			uc := &usecases.Usecases{
				OrderUsecases: order,
			}
			ctx := context.Background()
			handler := handler.New(ctx, uc)

			router := http.NewServeMux()
			router.HandleFunc("DELETE /remove", handler.Remove)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/remove", bytes.NewBufferString(tc.inputBody))

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
