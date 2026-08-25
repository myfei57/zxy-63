package coag

import "strconv"

func (d *Doser) DoseForFlow(value float64) float64 {
	return value * d.CurrentRatio()
}

func (d *Doser) DoseForTurbidity(turbidity float64) float64 {
	return turbidity * d.CurrentRatio()
}

func (d *Doser) DoseFromPersistedFlow() float64 {
	value, ok := d.flow.LoadFlow()
	if !ok {
		return 0
	}
	return d.DoseForFlow(value)
}

func (d *Doser) UpdateFlowAndDose(value float64) (float64, error) {
	if err := d.flow.PersistFlow(value); err != nil {
		return 0, err
	}
	dose := d.DoseFromPersistedFlow()
	if err := d.audit.Record("coagulant", strconv.FormatFloat(dose, 'f', 4, 64)); err != nil {
		return 0, err
	}
	return dose, nil
}
