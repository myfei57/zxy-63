package backwash

import "fmt"

func (c *Controller) Select(zone int) (string, error) {
	mapping := c.bank.Mapping()
	bedID, ok := mapping[zone]
	if !ok {
		return "", fmt.Errorf("no filter bed for zone %d", zone)
	}
	return bedID, nil
}
