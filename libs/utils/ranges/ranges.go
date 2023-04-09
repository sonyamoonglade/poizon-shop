package ranges

func IsBetween(x, lower, upper int) bool {
	return x > lower && x < upper
}

func IsBetweenInc(x, lower, upper int) bool {
	return x >= lower && x <= upper
}

func In(x int, arr []int) bool {
	for _, a := range arr {
		if x == a {
			return true
		}
	}
	return false
}
