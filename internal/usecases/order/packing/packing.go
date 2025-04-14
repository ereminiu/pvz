package packing

import (
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/usecases/order/packing/bag"
	"github.com/ereminiu/pvz/internal/usecases/order/packing/box"
	"github.com/ereminiu/pvz/internal/usecases/order/packing/film"
)

type Packing interface {
	CheckWeight(weight int) error
	GetPrice() int
}

func New(packingType string) (Packing, error) {
	switch packingType {
	case "film":
		return film.New(), nil

	case "bag":
		return bag.New(), nil

	case "box":
		return box.New(), nil

	default:
		return nil, myerrors.ErrInvalidOrderPackingType
	}
}
