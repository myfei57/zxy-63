package backwash

import "sort"

func (c *Controller) OrderRotation() ([]string, error) {
	beds := c.bank.Beds()
	sort.Slice(beds, func(i, j int) bool { return beds[i].Load < beds[j].Load })
	order := make([]string, 0, len(beds))
	for _, bed := range beds {
		order = append(order, bed.ID)
	}
	if len(order) > 0 {
		if err := c.bank.Rotate(order[0]); err != nil {
			return nil, err
		}
	}
	return order, nil
}
