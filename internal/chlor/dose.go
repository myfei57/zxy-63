package chlor

import "strconv"

func (d *Doser) DoseForResidual() float64 {
	return d.CurrentTarget()
}

func (d *Doser) ApplyResidual() error {
	dose := d.DoseForResidual()
	return d.audit.Record("chlorine", strconv.FormatFloat(dose, 'f', 4, 64))
}
