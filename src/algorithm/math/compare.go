package math

// 按要求比较两个数的大小

func compare(x, y int) bool {
	return (x-y)&(1<<31) == 0
}
