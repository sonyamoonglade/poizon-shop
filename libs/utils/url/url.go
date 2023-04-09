package url

import (
	"regexp"

	"functools"
)

var (
	r1 = regexp.MustCompile(`(?m)^https:\/\/dw4\.co\/.*`)
	r2 = regexp.MustCompile(`(?m)^https:\/\/qr\.1688\.com\/.*`)
	r3 = regexp.MustCompile(`(?m)^https:\/\/m\.tb\.cn\/.*`)
)

func IsValidDW4URL(url string) bool {
	return functools.Any(r1.MatchString(url), r2.MatchString(url), r3.MatchString(url))
}
