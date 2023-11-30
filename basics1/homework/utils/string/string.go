package string

func Populate(s string, n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += s
	}

	return res
}
