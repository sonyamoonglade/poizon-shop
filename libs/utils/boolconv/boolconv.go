package boolconv

import "strconv"

func Optimized(b bool) string {
	if b {
		return "1"
	}

	return "0"
}

func Full(b bool) string {
	return strconv.FormatBool(b)
}
