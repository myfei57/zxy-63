package clearwell

import "waterplant/internal/intake"

func (w *Well) AdjustLevel(target float64, inlet, outlet *intake.InletController) (float64, error) {
	start := w.Level()
	delta := target - start
	var mid float64
	var final float64
	if delta >= 0 {
		mid = inlet.Raise(start, delta)
		final = outlet.Lower(mid, delta)
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
