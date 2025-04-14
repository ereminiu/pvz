//go:build unit
// +build unit

package user_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/ereminiu/pvz/internal/handler"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/usecases"
	mock_usecases "github.com/ereminiu/pvz/internal/usecases/mocks"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestRefundOrders(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderIDs []int)

	type refundOrdersInput struct {
		UserID   int   `json:"user_id"`
		OrderIDs []int `json:"order_ids"`
	}

	testTable := []struct {
		name                 string
		input                refundOrdersInput
		call                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: refundOrdersInput{
				UserID:   19,
				OrderIDs: []int{26, 27},
			},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderIDs []int) {
				s.EXPECT().RefundOrders(ctx, userID, orderIDs).Return(nil)
			},
			expectedStatusCode:   http.StatusResetContent,
			expectedResponseBody: "orders have been refunded",
		},
		{
			name: "OK",
			input: refundOrdersInput{
				UserID:   19,
				OrderIDs: []int{89, 12},
			},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderIDs []int) {
				s.EXPECT().RefundOrders(ctx, userID, orderIDs).Return(myerrors.ErrInvalidOrderInput)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrInvalidOrderInput.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			user := mock_usecases.NewMockUserUsecases(c)
			tc.call(ctx, user, tc.input.UserID, tc.input.OrderIDs)

			uc := &usecases.Usecases{
				UserUsecases: user,
			}
			handler := handler.New(ctx, uc)

			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refund", bytes.NewBufferString(body))

			router := http.NewServeMux()
			router.HandleFunc("POST /refund", handler.RefundOrders)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestReturnOrders(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderIDs []int)

	type returnOrdersInput struct {
		UserID   int   `json:"user_id"`
		OrderIDs []int `json:"order_ids"`
	}

	testTable := []struct {
		name                 string
		input                returnOrdersInput
		call                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: returnOrdersInput{
				UserID:   19,
				OrderIDs: []int{26, 27},
			},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderIDs []int) {
				s.EXPECT().ReturnOrders(ctx, userID, orderIDs).Return(nil)
			},
			expectedStatusCode:   http.StatusResetContent,
			expectedResponseBody: "orders have been returned",
		},
		{
			name: "OK",
			input: returnOrdersInput{
				UserID:   19,
				OrderIDs: []int{89, 12},
			},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderIDs []int) {
				s.EXPECT().ReturnOrders(ctx, userID, orderIDs).Return(myerrors.ErrInvalidOrderInput)
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: myerrors.ErrInvalidOrderInput.Error() + "\n",
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			user := mock_usecases.NewMockUserUsecases(c)
			tc.call(ctx, user, tc.input.UserID, tc.input.OrderIDs)

			uc := &usecases.Usecases{
				UserUsecases: user,
			}
			handler := handler.New(ctx, uc)

			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/return", bytes.NewBufferString(body))

			router := http.NewServeMux()
			router.HandleFunc("POST /return", handler.ReturnOrders)

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
