package base

import (
	"math"
	"regexp"
)

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

var urlRegexp, _ = regexp.Compile(`^http[s]{0,1}://`)

func CompleteURL(u string) bool {
	return urlRegexp.MatchString(u)
}

var nameRegexp, _ = regexp.Compile(`_part_\d+`)

func ReplaceName(name string) string {
	return nameRegexp.ReplaceAllString(name, "")
}
