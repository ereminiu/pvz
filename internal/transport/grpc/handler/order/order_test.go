package order_test

import (
	"context"
	"testing"

	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/pb/api"
	"github.com/ereminiu/pvz/internal/transport/grpc/handler/order"
	mock_usecases "github.com/ereminiu/pvz/internal/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ereminiu/pvz/internal/entities"
)

func TestOrderServer_Create(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockOrderUsecases, order *entities.Order)

	testTable := []struct {
		name       string
		inputOrder *entities.Order
		call       mockBehavior
		Error      error
	}{
		{
			name: "OK",
			inputOrder: &entities.Order{
				UserID:      19,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       100,
				Packing:     "box",
				Extra:       false,
			},
			call: func(ctx context.Context, s *mock_usecases.MockOrderUsecases, order *entities.Order) {
				s.EXPECT().AddOrder(ctx, order).Return(nil)
			},
			Error: nil,
		},
		{
			name: "Invalid packing type",
			inputOrder: &entities.Order{
				UserID:      19,
				OrderID:     89,
				ExpireAfter: 2,
				Weight:      5,
				Price:       100,
				Packing:     "abobox",
				Extra:       false,
			},
			call: func(ctx context.Context, s *mock_usecases.MockOrderUsecases, order *entities.Order) {
				s.EXPECT().AddOrder(ctx, order).Return(myerrors.ErrInvalidOrderPackingType)
			},
			Error: myerrors.ErrInvalidOrderPackingType,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockOrderUsecases(c)
			tc.call(ctx, uc, tc.inputOrder)

			orderClient := order.New(uc)
			_, err := orderClient.Create(ctx, &api.AddOrderRequest{
				UserId:      int32(tc.inputOrder.UserID),
				OrderId:     int32(tc.inputOrder.OrderID),
				ExpireAfter: int32(tc.inputOrder.ExpireAfter),
				Weight:      int32(tc.inputOrder.Weight),
				Price:       int32(tc.inputOrder.Price),
				Packing:     tc.inputOrder.Packing,
				Extra:       tc.inputOrder.Extra,
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}

func TestOrderServer_Remove(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockOrderUsecases, orderID int)

	testTable := []struct {
		name         string
		inputOrderID int
		call         mockBehavior
		Error        error
	}{
		{
			name:         "OK",
			inputOrderID: 89,
			call: func(ctx context.Context, s *mock_usecases.MockOrderUsecases, orderID int) {
				s.EXPECT().RemoveOrder(ctx, orderID).Return(nil)
			},
			Error: nil,
		},
		{
			name:         "Invalid order_id",
			inputOrderID: 89,
			call: func(ctx context.Context, s *mock_usecases.MockOrderUsecases, orderID int) {
				s.EXPECT().RemoveOrder(ctx, orderID).Return(myerrors.ErrInvalidOrderInput)
			},
			Error: myerrors.ErrInvalidOrderInput,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockOrderUsecases(c)
			tc.call(ctx, uc, tc.inputOrderID)

			orderClient := order.New(uc)
			_, err := orderClient.Remove(ctx, &api.RemoveOrderRequest{
				OrderId: int32(tc.inputOrderID),
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}
