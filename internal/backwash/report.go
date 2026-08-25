package backwash

import (
	"fmt"
	"strings"

	"waterplant/internal/store"
)

type BackwashState struct {
	Commands []string `json:"commands"`
	Drains   []string `json:"drains"`
}

func (c *Controller) State() BackwashState {
	commands := store.LoadCommands(c.store, commandKey)
	var drains []string
	for _, key := range c.store.Keys() {
		if strings.HasPrefix(key, drainKeyPrefix) {
			drains = append(drains, strings.TrimPrefix(key, drainKeyPrefix))
		}
	}
	return BackwashState{Commands: commands, Drains: drains}
}

func (c *Controller) Describe() string {
	state := c.State()
	return fmt.Sprintf("backwash commands=%d drains=%d", len(state.Commands), len(state.Drains))
}

func (c *Controller) Eligible() []string {
	var eligible []string
	for _, bed := range c.bank.Beds() {
		if !bed.Closed {
			eligible = append(eligible, bed.ID)
		}
	}
	return eligible
}

func (c *Controller) DrainCount() int {
	return len(c.State().Drains)
}

func (c *Controller) PendingCount() int {
	return len(store.LoadCommands(c.store, commandKey))
}

func (c *Controller) LastDrain() (string, bool) {
	drains := c.State().Drains
	if len(drains) == 0 {
		return "", false
	}
	return drains[len(drains)-1], true
}

func (c *Controller) CommandList() []string {
	return store.LoadCommands(c.store, commandKey)
}
