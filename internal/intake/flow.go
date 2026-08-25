package intake

import (
	"strconv"

	"waterplant/internal/store"
)

const flowKey = "intake:flow"

type FlowRepository struct {
	store *store.Store
}

func NewFlowRepository(s *store.Store) *FlowRepository {
	return &FlowRepository{store: s}
}

func (r *FlowRepository) PersistFlow(value float64) error {
	return r.store.Put(flowKey, strconv.FormatFloat(value, 'f', 4, 64))
}

func (r *FlowRepository) LoadFlow() (float64, bool) {
	raw, ok := r.store.Get(flowKey)
	if !ok {
		return 0, false
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, false
	}
	return value, true
}
