package string

import (
	"strings"
)

func isSubstring(a, b string) bool {
	return strings.Index(a, b) != -1
}

func isRotatedString(a, b string) bool {
	al, bl := len(a), len(b)
	if al != bl {
		return false
	}

	var ai, bi int

	for bi < bl && a[ai] != b[bi] {
		bi++
	}

	// fmt.Println(a, b[bi:])

	for bi+ai < bl && a[ai] == b[bi+ai] {
		ai++
	}

	// fmt.Println(a[ai:], b[:bi])

	return a[ai:] == b[:bi]
}
