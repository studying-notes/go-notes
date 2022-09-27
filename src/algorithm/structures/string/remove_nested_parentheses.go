package string

import "bytes"

func RemoveNestedParentheses(s string) string {
	rs := []rune(s)

	if rs[0] != '(' || rs[len(rs)-1] != ')' {
		panic("invalid format")
	}

	var buf bytes.Buffer
	brackets := 0
	for i := range rs {
		if rs[i] == '(' {
			brackets++
		} else if rs[i] == ')' {
			brackets--
		} else {
			buf.WriteRune(rs[i])
		}
	}

	if brackets != 0 {
		panic("invalid format")
	}

	buf.WriteRune(')')

	return buf.String()
}
