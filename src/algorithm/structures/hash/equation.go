package hash

func FindEquation(list []int) ([2]int, [2]int) {
	// 两数和:两数 键值对
	kv := make(map[int][2]int)

	// 双重循环
	for idx, val := range list {
		for i := idx + 1; i < len(list); i++ {
			k := val + list[i]
			if v, ok := kv[k]; ok {
				return v, [2]int{val, list[i]}
			}
			kv[k] = [2]int{val, list[i]}
		}
	}
	return [2]int{0, 0}, [2]int{0, 0}
}
