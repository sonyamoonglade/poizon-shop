package functools

type mapFunc[A any, B any] func(a A, i int) B

func Map[A any, B any](f mapFunc[A, B], input []A) []B {
	result := make([]B, 0, len(input))

	for i := range input {
		result = append(result, f(input[i], i))
	}

	return result
}

type foreachFunc[A any] func(a A, i int) error

func ForEach[A any](f foreachFunc[A], input []A) error {
	for i, v := range input {
		if err := f(v, i); err != nil {
			return err
		}
	}
	return nil
}

func Reduce[IterType, AccType any](f func(acc AccType, el IterType) AccType, arr []IterType, acc AccType) AccType {
	for _, v := range arr {
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
