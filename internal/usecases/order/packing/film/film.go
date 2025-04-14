package film

const (
	price = 1
)

type Film struct{}

func New() *Film {
	return &Film{}
}

func (f *Film) CheckWeight(weight int) error {
	return nil
}

func (f *Film) GetPrice() int {
	return price
}
