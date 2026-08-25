package backwash

import "waterplant/internal/store"

const commandKey = "backwash:commands"

func (c *Controller) Enqueue(bedID string) error {
	return store.AppendCommand(c.store, commandKey, bedID)
}

func (c *Controller) Replay() ([]string, error) {
	commands := store.LoadCommands(c.store, commandKey)
	for _, bedID := range commands {
		if err := c.Start(bedID); err != nil {
			return commands, err
		}
	}
	return commands, nil
}

func (c *Controller) Recover() error {
	return nil
}
