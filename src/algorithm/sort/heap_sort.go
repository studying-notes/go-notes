package sort

import (
	. "algorithm/structures/heap"
)

func MaxHeapSort(array []int) {
	length := len(array)

	// 建立最大堆
	ConvertArrayToMaxHeap(array)

	for i := length - 1; i >= 0; i-- {
		// 顶部元素一定是最大的，所以每次都将其放到最后
		array[0], array[i] = array[i], array[0]
		// 交换堆顶元素和最后一个元素，然后调整剩余堆
		// 在调整时排除已经排序的元素，所以将 i 作为长度传入
		AdjustToMaxHeap(array, 0, i)
	}
}
