package bag

import (
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/pkg/errors"
)

const (
	maxWeight = 10
	price     = 5
)

type Bag struct{}

func New() *Bag {
	return &Bag{}
}

func (f *Bag) CheckWeight(weight int) error {
	if weight >= maxWeight {
		return errors.Wrapf(myerrors.ErrInvalidOrderWeight, "weight should be less than %d", maxWeight)
	}
	return nil
}

func (f *Bag) GetPrice() int {
	return price
}
