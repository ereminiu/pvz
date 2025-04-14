package filter

import (
	"github.com/ereminiu/pvz/internal/pkg/auditlog/models"
)

type Empty struct {
}

func NewEmpty() *Empty {
	return &Empty{}
}

func (empty *Empty) Check(event models.Log) bool {
	return true
}

type Action struct {
	shouldBe string
}

func NewAction(shouldBe string) *Action {
	return &Action{
		shouldBe: shouldBe,
	}
}

func (action *Action) Check(event models.Log) bool {
	return action.shouldBe == "" || event.Action == action.shouldBe
}
