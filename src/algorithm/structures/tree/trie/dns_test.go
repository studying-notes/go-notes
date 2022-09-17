package trie

import "fmt"

func ExampleDNSSearch() {
	ips := []string{"248.116.89.121", "89.105.17.198",
		"69.204.3.67", "188.127.67.5", "73.255.192.234"}
	domains := []string{"www.samsung.com", "www.samsung.net",
		"www.baidu.cn", "google.com", "google.com"}

	c := NewDNSCache()
	for i, v := range ips {
		c.Add(v, domains[i])
	}

	ip := ips[1]

	fmt.Println(c.Get(ip))

	// Output:
	// www.samsung.net
}
