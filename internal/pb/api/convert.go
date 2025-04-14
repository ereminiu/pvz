package api

import (
	"github.com/ereminiu/pvz/internal/entities"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func ToOrderList(orders []*entities.Order) []*OrderInfo {
	if orders == nil {
		return nil
	}

	orderInfos := make([]*OrderInfo, 0, len(orders))
	for _, order := range orders {
		orderInfos = append(orderInfos, ToOrderInfo(order))
	}

	return orderInfos
}

func ToOrderInfo(order *entities.Order) *OrderInfo {
	return &OrderInfo{
		OrderId:     int32(order.OrderID),
		UserId:      int32(order.UserID),
		Weight:      int32(order.Weight),
		Price:       int32(order.Price),
		Packing:     order.Packing,
		Extra:       order.Extra,
		Status:      order.Status,
		ExpireAfter: int32(order.ExpireAfter),
		ExpireAt:    timestamppb.New(order.ExpireAt),
		UpdateAt:    timestamppb.New(order.UpdatedAt),
		CreatedAt:   timestamppb.New(order.CreatedAt),
	}
}
