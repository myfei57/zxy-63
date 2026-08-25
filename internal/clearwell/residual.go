package clearwell

func (w *Well) UpdateResidualDemand(demand float64) error {
	target := 0.3 + demand*0.7
	return w.SetResidualTarget(target)
}
