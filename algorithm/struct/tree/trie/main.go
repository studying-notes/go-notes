package main

import (
	"fmt"
	"github.com/go-ego/cedar"
)

func main() {
	trie := cedar.New()

	printIdKeyValue := func(id int) {
		key, _ := trie.Key(id)
		value, _ := trie.Value(id)
		fmt.Printf("%d\t%s:%v\n", id, key, value)
	}

	_ = trie.Insert([]byte("How many"), 0)
	_ = trie.Insert([]byte("How many loved"), 1)
	_ = trie.Insert([]byte("How many loved your moments"), 2)
	_ = trie.Insert([]byte("How many loved your moments of glad grace"), 3)

	_ = trie.Insert([]byte("姑苏"), 4)
	_ = trie.Insert([]byte("姑苏城外"), 5)
	_ = trie.Insert([]byte("姑苏城外寒山寺"), 6)

	value, _ := trie.Get([]byte("How many loved your moments of glad grace"))
	fmt.Println(value)

	id, _ := trie.Jump([]byte("How many loved your moments"), 0)
	printIdKeyValue(id)

	// 输入字符串开始部分匹配 key
	fmt.Println("\nPrefixMatch\nid\tkey:value")
	for _, id := range trie.PrefixMatch([]byte("How many loved your moments of glad grace"), 0) {
		printIdKeyValue(id)
	}

	// key 开始部分匹配输入字符串
	fmt.Println("\nPrefixPredict\nid\tkey:value")
	for _, id := range trie.PrefixPredict([]byte("姑苏"), 0) {
		printIdKeyValue(id)
	}
}
