package quota

type Quota struct {
	Chemical string
	Limit    float64
}

func CheckQuota(chemical string, used, limit float64) (float64, bool) {
	if limit <= 0 {
		return used, true
	}
	remaining := limit - used
	return remaining, remaining >= 0
}
