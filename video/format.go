package video

import (
	"math"
	"regexp"
)

func unitReturn(value float64) (float64, string) {
	switch {
	case value > 1024:
		// kb
		return value / 1024, "kb"
	case value > math.Pow(1024, 2):
		// mb
		return value / math.Pow(1024, 2), "mb"
	case value > math.Pow(1024, 3):
		// gb
		return value / math.Pow(1024, 3), "gb"
	default:
		// byte
		return value, "byte"
	}
}

var urlRegexp, _ = regexp.Compile(`^http[s]{0,1}://`)

func CompleteURL(u string) bool {
	return urlRegexp.MatchString(u)
}
