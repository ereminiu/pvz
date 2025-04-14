package collections

import (
	"time"

	"github.com/ereminiu/pvz/internal/entities"
)

var (
	Order1 = &entities.Order{
		UserID:   19,
		OrderID:  89,
		ExpireAt: time.Now().AddDate(0, 0, 3),
		Weight:   5,
		Price:    100,
		Packing:  "box",
		Extra:    true,
		Status:   "delivered",
	}
	Order2 = &entities.Order{
		UserID:   19,
		OrderID:  86,
		ExpireAt: time.Now().AddDate(0, 0, 2),
		Weight:   5,
		Price:    100,
		Packing:  "film",
		Extra:    false,
		Status:   "delivered",
	}
	Order3 = &entities.Order{
		UserID:   19,
		OrderID:  27,
		ExpireAt: time.Now().AddDate(0, 0, 7),
		Weight:   7,
		Price:    101,
		Packing:  "bag",
		Extra:    true,
		Status:   "delivered",
	}

	GivenOrder1 = &entities.Order{
		UserID:   26,
		OrderID:  94,
		ExpireAt: time.Now().AddDate(0, 0, 7),
		Weight:   7,
		Price:    101,
		Packing:  "bag",
		Extra:    true,
		Status:   "given",
	}
	GivenOrder2 = &entities.Order{
		UserID:   26,
		OrderID:  93,
		ExpireAt: time.Now().AddDate(0, 0, 7),
		Weight:   3,
		Price:    128,
		Packing:  "box",
		Extra:    true,
		Status:   "given",
	}

	RefundOrder1 = &entities.Order{
		UserID:   19,
		OrderID:  12,
		ExpireAt: time.Now().AddDate(0, 0, 2),
		Weight:   8,
		Price:    200,
		Packing:  "box",
		Extra:    false,
		Status:   "refund",
	}
	RefundOrder2 = &entities.Order{
		UserID:   19,
		OrderID:  112,
		ExpireAt: time.Now().AddDate(0, 0, 2),
		Weight:   8,
		Price:    230,
		Packing:  "film",
		Extra:    false,
		Status:   "refund",
	}
	RefundOrder3 = &entities.Order{
		UserID:   19,
		OrderID:  1112,
		ExpireAt: time.Now().AddDate(0, 0, 2),
		Weight:   8,
		Price:    200,
		Packing:  "film",
		Extra:    false,
		Status:   "refund",
	}
)

var (
	DeliveredOrders = []*entities.Order{Order1, Order2, Order3}
	GivenOrders     = []*entities.Order{GivenOrder1, GivenOrder2}
	RefundOrders    = []*entities.Order{RefundOrder1, RefundOrder2, RefundOrder3}
)
