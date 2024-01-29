package slice

func Filter[T any](ss []T, test func(T) bool) []T {
	res := make([]T, 0, len(ss))

	for _, s := range ss {
		if test(s) {
			res = append(res, s)
		}
	}

	return res
}

func ContainsValue[T comparable](ss []*T, v T) bool {
	for _, s := range ss {
		if *s == v {
			return true
		}
	}

	return false
}

func Map[T any, V any](s []T, f func(T) V) []V {
	res := make([]V, len(s))

	for i, e := range s {
		res[i] = f(e)
	}

	return res
}
