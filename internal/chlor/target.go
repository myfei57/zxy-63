package chlor

func (d *Doser) CurrentTarget() float64 {
	return d.well.ResidualTarget()
}
