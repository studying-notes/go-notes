package heap

import "fmt"

// AdjustToMaxHeap 将指定节点调整为最大堆
// parent 指定节点的索引
// length 指定节点的数组长度
func AdjustToMaxHeap(array []int, parent, length int) {
	for parent*2+1 < length {
		// 左子节点存在
		child := parent*2 + 1

		if child+1 < length && array[child+1] > array[child] {
			// 存在右子节点，且右子节点的值更大
			child++
		}

		if array[child] > array[parent] {
			// 子节点的值大于父节点
			array[parent], array[child] = array[child], array[parent]
			// 父子节点交换了数据，所以必须继续调整该子节点
			parent = child
		} else {
			// 说明子节点已经有序，中断循环
			break
		}
	}
}

// ConvertArrayToMaxHeap 将数组转换为最大堆
func ConvertArrayToMaxHeap(array []int) {
	length := len(array)

	// 建立最大堆
	for i := length/2 - 1; i >= 0; i-- {
		// 从最底部的父节点开始遍历
		AdjustToMaxHeap(array, i, length)
		fmt.Printf("array[%d]=%d\n", i, array[i])
		fmt.Println(array)
	}
}

func PrintHeap(array []int) {

}
