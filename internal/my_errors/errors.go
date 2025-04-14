package myerrors

import "errors"

var (
	ErrOrderAlreadyCreated     = errors.New("order already created")
	ErrOrderAlreadyRemoved     = errors.New("order already removed")
	ErrInvalidOrderInput       = errors.New("invalid order input")
	ErrInvalidOrderPackingType = errors.New("invalid packing type")
	ErrInvalidOrderWeight      = errors.New("invalid weight for this type of packign")

	ErrIncorrectPassword = errors.New("incorrect password")

	ErrDuringWritingResponse = errors.New("error during writing response")
)
