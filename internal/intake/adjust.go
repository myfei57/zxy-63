package intake

type InletController struct{}

func NewInletController() *InletController {
	return &InletController{}
}

func (c *InletController) Raise(level, amount float64) float64 {
	next := level + amount
	if next < 0 {
		return 0
	}
	return next
}

func (c *InletController) Lower(level, amount float64) float64 {
	next := level - amount
	if next < 0 {
		return 0
	}
	return next
}
