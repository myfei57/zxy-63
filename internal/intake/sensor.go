package intake

type Sensor struct {
	FlowValue float64 `json:"flow"`
	Turbidity float64 `json:"turbidity"`
}
