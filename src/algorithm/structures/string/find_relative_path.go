package string

func findRelativePath(a, b string) string {
	al, bl := len(a), len(b)

	var i, j int
	var s = "../"

	for i < al && i < bl && a[i] == b[i] {
		i++
	}

	for j = i; j < al; j++ {
		if a[j] == '/' {
			s += "../"
		}
	}

	return s + b[i:]
}
