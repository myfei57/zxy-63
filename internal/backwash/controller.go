package backwash

import (
	"waterplant/internal/filter"
	"waterplant/internal/store"
)

type Controller struct {
	bank    *filter.Bank
	store   *store.Store
	mapping map[int]string
}

func NewController(b *filter.Bank, s *store.Store) *Controller {
	return &Controller{bank: b, store: s, mapping: b.Mapping()}
}
