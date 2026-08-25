package flow

type Meter struct {
	Serial string
	Factor float64
}

func CalibrateMeter(serial string, factor float64) Meter {
	return Meter{Serial: serial, Factor: factor}
}
