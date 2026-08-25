package intake

func Mix(samples []float64) float64 {
	if len(samples) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range samples {
		total += value
	}
	return total / float64(len(samples))
}
