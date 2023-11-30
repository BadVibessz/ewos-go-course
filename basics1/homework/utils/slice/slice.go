package slice

// todo: generic
func MapStringToInt(s []string, f func(string) int) []int {

	res := make([]int, len(s))
	for i := range s {
		res[i] = f(s[i])
	}
	return res
}
