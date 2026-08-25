package clearwell

import "waterplant/internal/intake"

func (w *Well) AdjustLevel(target float64, inlet, outlet *intake.InletController) (float64, error) {
	start := w.Level()
	delta := target - start
	var mid float64
	var final float64
	if delta >= 0 {
		mid = outlet.Lower(start, delta)
		final = inlet.Raise(mid, delta)
	} else {
		amount := -delta
		mid = inlet.Raise(start, amount)
		final = outlet.Lower(mid, amount)
	}
	if err := w.SetLevel(final); err != nil {
		return 0, err
	}
	return min(start, mid, final), nil
}
