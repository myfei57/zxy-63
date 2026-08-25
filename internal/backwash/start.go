package backwash

import (
	"strconv"
)

const drainKeyPrefix = "backwash:drain:"
const spillKey = "backwash:spill"

func (c *Controller) Start(bedID string) error {
	if err := c.drain(bedID); err != nil {
		return err
	}
	return c.bank.Close(bedID)
}

func (c *Controller) drain(bedID string) error {
	closed, err := c.bank.IsClosed(bedID)
	if err != nil {
		return err
	}
	if !closed {
		if err := c.store.Put(spillKey, "true"); err != nil {
			return err
		}
	}
	return c.store.Put(drainKeyPrefix+bedID, strconv.FormatInt(nowUnix(), 10))
}
