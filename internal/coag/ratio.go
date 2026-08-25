package coag

func (d *Doser) CurrentRatio() float64 {
	return d.calib.Current()
}
