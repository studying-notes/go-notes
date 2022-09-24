package array

// IsEqualArray 比较两个切片是否相等 reflect.DeepEqual 性能极差
func IsEqualArray(a, b []int) bool {
	if len(a) != len(b) || (a == nil) != (b == nil) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
