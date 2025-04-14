package order

import "github.com/ereminiu/pvz/internal/entities"

type Packing interface {
	CheckWeight(weight int) error
	GetPrice() int
}

type Roller struct {
	order *entities.Order
}

func (r *Roller) Roll(packing Packing) error {
	if err := packing.CheckWeight(r.order.Weight); err != nil {
		return err
	}

	r.order.Price += packing.GetPrice()

	return nil
}
