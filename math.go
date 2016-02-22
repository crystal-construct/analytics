package analytics

import "math"

func round(f float64) float64 {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return float64(int(f + math.Copysign(0.5, f)))
}
