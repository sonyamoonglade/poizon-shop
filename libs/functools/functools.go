package functools

type mapFunc[A any, B any] func(a A, i int) B

func Map[A any, B any](f mapFunc[A, B], input []A) []B {
	result := make([]B, 0, len(input))

	for i := range input {
		result = append(result, f(input[i], i))
	}

	return result
}

func Reduce[A, B any](f func(B, A) B, s []A, initValue B) B {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

func Any(bits ...bool) bool {
	for i := 0; i < len(bits); i++ {
		if bits[i] {
			return true
		}
	}
	return false
}
