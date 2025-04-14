package user_test

import (
	"context"
	"testing"

	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pb/api"
	"github.com/ereminiu/pvz/internal/transport/grpc/handler/user"
	mock_usecases "github.com/ereminiu/pvz/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestServer_Refund(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderID []int)

	testTable := []struct {
		name         string
		inputUserID  int
		inputOrderID []int
		call         mockBehavior
		Error        error
	}{
		{
			name:         "OK",
			inputUserID:  19,
			inputOrderID: []int{89, 12},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderID []int) {
				s.EXPECT().RefundOrders(ctx, userID, orderID).Return(nil)
			},
			Error: nil,
		},
		{
			name:         "Invalid order_id",
			inputUserID:  19,
			inputOrderID: []int{89, 12},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderID []int) {
				s.EXPECT().RefundOrders(ctx, userID, orderID).Return(myerrors.ErrInvalidOrderInput)
			},
			Error: myerrors.ErrInvalidOrderInput,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockUserUsecases(c)
			tc.call(ctx, uc, tc.inputUserID, tc.inputOrderID)

			userClient := user.New(uc)
			_, err := userClient.Refund(ctx, &api.RefundOrderRequest{
				UserId:  int32(tc.inputUserID),
				OrderId: convertOrderID(tc.inputOrderID),
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}

func TestServer_Return(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderID []int)

	testTable := []struct {
		name         string
		inputUserID  int
		inputOrderID []int
		call         mockBehavior
		Error        error
	}{
		{
			name:         "OK",
			inputUserID:  19,
			inputOrderID: []int{89, 12},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderID []int) {
				s.EXPECT().ReturnOrders(ctx, userID, orderID).Return(nil)
			},
			Error: nil,
		},
		{
			name:         "Invalid order_id",
			inputUserID:  19,
			inputOrderID: []int{89, 12},
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID int, orderID []int) {
				s.EXPECT().ReturnOrders(ctx, userID, orderID).Return(myerrors.ErrInvalidOrderInput)
			},
			Error: myerrors.ErrInvalidOrderInput,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockUserUsecases(c)
			tc.call(ctx, uc, tc.inputUserID, tc.inputOrderID)

			userClient := user.New(uc)
			_, err := userClient.Return(ctx, &api.ReturnOrderRequest{
				UserId:  int32(tc.inputUserID),
				OrderId: convertOrderID(tc.inputOrderID),
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}

func TestServer_List(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID, lastN int, located bool, pattern map[string]string)

	testTable := []struct {
		name    string
		userID  int
		lastN   int
		located bool
		pattern map[string]string
		call    mockBehavior
		Error   error
	}{
		{
			name:    "OK",
			userID:  19,
			lastN:   5,
			located: false,
			pattern: nil,
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID, lastN int, located bool, pattern map[string]string) {
				s.EXPECT().GetList(ctx, userID, lastN, located, pattern).Return(nil, nil)
			},
			Error: nil,
		},
		{
			name:    "Invalid user_id",
			userID:  19,
			lastN:   5,
			located: false,
			pattern: nil,
			call: func(ctx context.Context, s *mock_usecases.MockUserUsecases, userID, lastN int, located bool, pattern map[string]string) {
				s.EXPECT().GetList(ctx, userID, lastN, located, pattern).Return(nil, myerrors.ErrInvalidOrderInput)
			},
			Error: myerrors.ErrInvalidOrderInput,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockUserUsecases(c)
			tc.call(ctx, uc, tc.userID, tc.lastN, tc.located, tc.pattern)

			userClient := user.New(uc)
			_, err := userClient.List(ctx, &api.ListRequest{
				UserId:  int32(tc.userID),
				LastN:   int32(tc.lastN),
				Located: tc.located,
				Pattern: tc.pattern,
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}

func convertOrderID(orderID []int) []int32 {
	res := make([]int32, 0, len(orderID))
	for _, x := range orderID {
		res = append(res, int32(x))
	}

	return res
}
