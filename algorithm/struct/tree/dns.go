package main

import "fmt"

// Trie 树定义
type TrieNode struct {
	IsLeaf bool
	Url    string
	Child  []*TrieNode
}

func NewTrieNode(count int) *TrieNode {
	return &TrieNode{
		IsLeaf: false,
		Url:    "",
		Child:  make([]*TrieNode, count),
	}
}

var CharCount = 11

type DNSCache struct {
	root *TrieNode
}

func (p *DNSCache) getIndexFromRune(r rune) int {
	if r == '.' {
		return 10
	} else {
		return int(r) - '0'
	}
}

func (p *DNSCache) getRuneFromIndex(i int) rune {
	if i == 10 {
		return '.'
	} else {
		return rune('0' + i)
	}
}

// 把一个 IP 地址和相应的 URL 添加到 Trie 树中，最后一个结点是 URL
func (p *DNSCache) Insert(ip, url string) {
	root := p.root
	for _, v := range []rune(ip) {
		// 根据当前遍历到的 IP 中的字符，找出子结点的索引
		index := p.getIndexFromRune(v)
		if root.Child[index] == nil {
			root.Child[index] = NewTrieNode(CharCount)
		}
		// 移动到子结点
		root = root.Child[index]
	}
	// 在叶子结点中存储 IP 地址对应的 URL
	root.IsLeaf = true
	root.Url = url
}

// 通过 IP 地址找到对应的 URL
func (p *DNSCache) SearchDNSCache(ip string) string {
	root := p.root
	for _, v := range []rune(ip) {
		index := p.getIndexFromRune(v)
		if root.Child[index] == nil {
			return ""
		}
		root = root.Child[index]
	}
	// 返回找到的 URL
	if root != nil && root.IsLeaf {
		return root.Url
	}
	return ""
}

func NewDNSCache() *DNSCache {
	return &DNSCache{root: NewTrieNode(CharCount)}
}

func main() {
	ipAddrs := []string{"248.116.89.121", "89.105.17.198",
		"69.204.3.67", "188.127.67.5", "73.255.192.234"}
	urls := []string{"www.samsung.com", "www.samsung.net",
		"www.baidu.cn", "google.com", "google.com"}
	c := NewDNSCache()
	for i, v := range ipAddrs {
		c.Insert(v, urls[i])
	}
	ip := ipAddrs[1]
	fmt.Println(c.SearchDNSCache(ip))
}
