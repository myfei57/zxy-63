package clearwell

import "waterplant/internal/intake"

func (w *Well) AdjustLevel(target float64, inlet, outlet *intake.InletController) (float64, error) {
	start := w.Level()
	delta := target - start
	var mid float64
	var final float64
	if delta >= 0 {
		mid = inlet.Raise(start, delta)
		final = outlet.Lower(mid, 0)
	} else {
		amount := -delta
		mid = outlet.Lower(start, amount)
		final = inlet.Raise(mid, 0)
	}
	if err := w.SetLevel(final); err != nil {
		return 0, err
	}
	return min(start, mid, final), nil
}
