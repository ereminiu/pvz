package box

import (
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/pkg/errors"
)

const (
	maxWeight = 30
	price     = 20
)

type Box struct{}

func New() *Box {
	return &Box{}
}

func (f *Box) CheckWeight(weight int) error {
	if weight >= maxWeight {
		return errors.Wrapf(myerrors.ErrInvalidOrderWeight, "weight should be less than %d", maxWeight)
	}
	return nil
}

func (f *Box) GetPrice() int {
	return price
}
