package numeric

var digits = make(map[rune]struct{})

func init() {
	for _, d := range "0123456789" {
		digits[d] = struct{}{}
	}
}

func IsDigit(r rune) bool {
	_, ok := digits[r]
	return ok
}

func AllAreDigits(s string) bool {
	for _, r := range s {
		if !IsDigit(r) {
			return false
		}
	}
	return true
}
