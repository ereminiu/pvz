package entities

import "time"

type Order struct {
	OrderID     int       `json:"order_id" db:"order_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	ExpireAfter int       `json:"expire_after" db:"-"`
	Weight      int       `json:"weight" db:"weight"`
	Price       int       `json:"price" db:"price"`
	Packing     string    `json:"packing" db:"packing"`
	Extra       bool      `json:"extra" db:"extra"`
	Status      string    `json:"status" db:"status"`
	ExpireAt    time.Time `json:"expire_at" db:"expire_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
