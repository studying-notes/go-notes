package array

// Range 可生成指定范围的自然数切片
func Range(args ...uint) []uint {
	var (
		start, end uint
		step       uint = 1
	)

	if len(args) == 0 {
		return []uint{}
	} else if len(args) == 1 {
		end = args[0]
	} else if len(args) == 2 {
		start, end = args[0], args[1]
	} else if len(args) > 2 {
		start, end, step = args[0], args[1], args[2]
	}

	if step == 0 ||
		start == end ||
		(step < 0 && start < end) ||
		(step > 0 && start > end) {
		return []uint{}
	}

	s := make([]uint, 0, (end-start)/step+1)

	for start < end {
		s = append(s, start)
		start += step
	}

	return s
}
