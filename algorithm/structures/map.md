---
date: 2020-10-12T17:08:42+08:00  # 创建日期
author: "Rustle Karl"  # 作者

title: "数据结构与算法之哈希"  # 文章标题
url:  "posts/go/algorithm/structures/map"  # 设置网页永久链接
tags: [ "algorithm", "go" ]  # 标签
categories: [ "Go 数据结构与算法"]  # 系列

weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

## 从给定的车票中找出旅程路线

给定一趟旅途旅程中所有的车票信息，根据这个车票信息找出这趟旅程的路线。

例如，给定车票（“西安”到“成都”），（“北京”到“上海”），（“大连”到“西安”），（“上海”到“大连”）。

那么可以得到旅程路线为：北京->上海， 上海->大连， 大连->西安， 西安->成都。假定给定的车票不会有环，也就是说有一个城市只作为终点而不会作为起点。

一般而言可以使用拓扑排序进行解答。根据车票信息构建一个图，然后找出这张图的拓扑排序序列，这个序列就是旅程的路线。

### 构建 Map 和逆 Map

```go
func RevTickets(tickets map[int]int) map[int]int {
	revTickets := make(map[int]int)
	for k, v := range tickets {
		revTickets[v] = k
	}
	return revTickets
}

func PrintTravel(tickets map[int]int) {
	revTickets := RevTickets(tickets)
	var next int
	for k := range tickets {
		_, ok := revTickets[k]
		if !ok {
			next = k
			break
		}
	}
	fmt.Print(next, " ")
	for range tickets {
		next = tickets[next]
		fmt.Print(next, " ")
	}
}
```

## 从数组中找出满足 `a+b=c+d` 的两个数对

给定一个数组，找出数组中是否有两个数对 (a, b) 和 (c, d)，使得 a+b=c+d，其中，a、b、c 和 d 是不同的元素。如果有多个答案，打印任意一个即可。例如给定数组：{3, 4, 7, 10, 20, 9, 8}，可以找到两个数对 (3, 8) 和 (4, 7)，使得 3+8=4+7。

```go
func main() {
	list := []int{3, 4, 7, 10, 2, 9, 8}
	fmt.Println(FindEquation(list))
}

func FindEquation(list []int) ([2]int, [2]int) {
	kv := make(map[int][2]int)
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
```
