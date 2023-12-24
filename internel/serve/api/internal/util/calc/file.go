package calc

import "math"

func UnitReturn(value float64) (float64, string) {
	switch {
	case value > math.Pow(1024, 3):
		// gb
		return value / math.Pow(1024, 3), "gb"
	case value > math.Pow(1024, 2):
		// mb
		return value / math.Pow(1024, 2), "mb"
	case value > 1024:
		// kb
		return value / 1024, "kb"
	default:
		// byte
		return value, "byte"
	}
}
